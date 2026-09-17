// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package main

import (
	"context"
	"flag"
	"os"

	"TravelHome/app/mqueue/cmd/job/internal/config"
	"TravelHome/app/mqueue/cmd/job/internal/logic"

	"TravelHome/app/mqueue/cmd/job/internal/svc"

	"github.com/zeromicro/go-zero/core/conf"
)

var configFile = flag.String("f", "etc/job-api.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	// log、prometheus、trace、metricsUrl
	if err := c.SetUp(); err != nil {
		panic(err)
	}

	svcContext := svc.NewServiceContext(c)
	ctx := context.Background()

	cornJob := logic.NewCronJob(ctx, svcContext)

	mux := cornJob.Register()

	if err := svcContext.AsynqServer.Run(mux); err != nil {
		// logx.WithContext(ctx).Errorf("!!!CronJobErr!!! run err:%+v", err)
		os.Exit(1)
	}

}
