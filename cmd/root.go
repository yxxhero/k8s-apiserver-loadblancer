package cmd

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"k8s.io/client-go/tools/leaderelection"

	"github.com/yxxhero/k8s-apiserver-loadblancer/pkg/config"
	"github.com/yxxhero/k8s-apiserver-loadblancer/pkg/k8s"
	"github.com/yxxhero/k8s-apiserver-loadblancer/pkg/lock"
	"github.com/yxxhero/k8s-apiserver-loadblancer/pkg/mirror"
)

func NewRootCmd() *cobra.Command {
	c := config.NewConfig()
	// rootCmd represents the base command when called without any subcommands
	rootCmd := &cobra.Command{
		Use:   "k8s-apiserver-loadblancer",
		Short: "it is a tool to create a loadbalancer for k8s apiserver",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, cancel := context.WithCancel(context.Background())
			sigCh := make(chan os.Signal, 1)
			stopCh := make(chan struct{})
			log.Debug().Msg("register signal handler")
			signal.Notify(sigCh, syscall.SIGTERM, syscall.SIGINT)
			go func() {
				<-sigCh
				log.Debug().Msg("received signal, stopping")
				close(stopCh)
				cancel()
				log.Info().Msg("received signal, exit...")
			}()
			c.StopCh = stopCh
			// set config ID
			if c.ID == "" {
				c.ID = os.Getenv("POD_NAME")
			}

			// verify the configuration
			if err := c.Verify(); err != nil {
				return err
			}
			k8sClient, err := k8s.NewClient(c.Kubeconfig)
			if err != nil {
				return err
			}

			resourceLock := lock.NewResourceLock(k8sClient, c)
			leaderelectionConfig := leaderelection.LeaderElectionConfig{
				Lock:            resourceLock,
				LeaseDuration:   60 * time.Second,
				RenewDeadline:   30 * time.Second,
				RetryPeriod:     10 * time.Second,
				ReleaseOnCancel: true,
				Callbacks: leaderelection.LeaderCallbacks{
					OnStartedLeading: func(ctx context.Context) {
						log.Info().Msg("started leading")
						if err := mirror.Run(ctx, c, k8sClient); err != nil {
							log.Error().Err(err).Msg("failed to run mirror")
						}
						log.Info().Msg("mirror stopped")
					},
					OnStoppedLeading: func() {
						log.Info().Msg("stopped leading")
					},
					OnNewLeader: func(identity string) {
						log.Info().Str("leader", identity).Msg("new leader")
					},
				},
			}
			leaderelection.RunOrDie(ctx, leaderelectionConfig)
			return nil
		},
	}
	rootCmd.Flags().StringVar(&c.Kubeconfig, "kubeconfig", "", "absolute path to the kubeconfig file")
	rootCmd.Flags().StringVar(&c.ServiceType, "service-type", "LoadBalancer", "service type")
	rootCmd.Flags().StringVar(&c.ID, "id", "", "unique id for the leader election")
	rootCmd.Flags().StringVar(&c.ServiceName, "service-name", "k8s-apiserver", "service name")
	rootCmd.Flags().StringVar(&c.ServiceNamespace, "service-namespace", "default", "service namespace")
	return rootCmd
}
