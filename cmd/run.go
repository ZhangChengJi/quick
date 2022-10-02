package cmd

import (
	"fmt"
	"github.com/lithammer/dedent"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"quick/conf"
	"quick/pkg/log"
	"syscall"
)

func init() {
	conf.InitConf()
	zap.ReplaceGlobals(log.New())
}

func Run() error {
	cmd := &cobra.Command{
		Use:   "giot 设备接入平台",
		Short: "giot",
		Long: dedent.Dedent(`
				┌──────────────────────────────────────────────────────────┐
			    │ FlYING                                                   │
			    │ Cloud Native Distributed Configuration Center             │
			    │                                                          │
			    │ Please give us feedback at:                              │
			    │ https://github.com/ZhangChengJi/flyingv2/issues           │
			    └──────────────────────────────────────────────────────────┘
		`),
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE: func(cmd *cobra.Command, args []string) error {
			conf.InitConf()
			zap.ReplaceGlobals(log.New())
			return nil
		},
	}
	cmd.PersistentFlags().StringVarP(&conf.ConfigFile, "config", "c", "", "config file")
	cmd.PersistentFlags().StringVarP(&conf.WorkDir, "work-dir", "p", ".", "current work directory")
	cmd.Execute()
	s := NewServer()
	errSig := make(chan error, 5)
	s.Start(errSig)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	select {
	case sig := <-quit:
		log.Sugar.Infof("The Manager API server receive %s and start shutting down", sig.String())
		s.Stop()
		log.Sugar.Info("See you next time!")
	case err := <-errSig:
		log.Sugar.Errorf("The Manager API server start failed: %s", zap.Error(err))
		return err
	}
	return nil
}

type server struct {
}

func NewServer() *server {
	return &server{}
}
func (s *server) Start(er chan error) {
	//err := s.init()
	//if err != nil {
	//	er <- err
	//	return
	//}
	s.printInfo()
}
func (s *server) Stop() {
	//s.shutdownServer(nil)

}
func (s *server) printInfo() {
	fmt.Fprint(os.Stdout, "quick running successfully!\n\n")
}
