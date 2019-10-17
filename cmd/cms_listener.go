package cmd

import (
	"context"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"os"
	"os/signal"
	"playhead/queues"
	"sync"
)

func listenCmsQueue(ctx context.Context, q *queues.Queue) {
	qctx := q.NewContext()
	q.Context = qctx
	q.StartCmsListener(qctx)
	done := make(chan struct{})
	go func() {
		<-ctx.Done()
		close(done)
	}()
	<-done
}

var listenCmsQ = &cobra.Command{
	Use:   "cms_listener",
	Short: "Start the CMS queue listener",
	Long: "Listens for messages from the CMS.  " +
		"These messages will dictate the relationship between series and episodes",
	RunE: func(cmd *cobra.Command, args []string) error {
		q, err := queues.New()
		if err != nil {
			return err
		}
		defer q.Close()

		ctx, cancel := context.WithCancel(context.Background())

		go func() {
			ch := make(chan os.Signal, 1)
			signal.Notify(ch, os.Interrupt)
			<-ch
			logrus.Info("signal caught. shutting down... I should really push things back into the queue but I DONT")
			cancel()
		}()

		var wg sync.WaitGroup

		wg.Add(1)
		go func() {
			defer wg.Done()
			defer cancel()
			listenCmsQueue(ctx, q)
		}()

		wg.Wait()
		return nil
	},
}

func init() {
	rootCmd.AddCommand(listenCmsQ)
}
