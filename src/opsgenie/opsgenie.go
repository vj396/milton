package opsgenie

import (
	"context"
	"fmt"
	"net/http"
	"sync/atomic"
	"time"

	"github.com/opsgenie/opsgenie-go-sdk-v2/client"
	"github.com/opsgenie/opsgenie-go-sdk-v2/team"
	"github.com/vj396/milton/src/types"
)

var (
	retries int32 = 0
	config        = &client.Config{}
)

func Set(conf *types.Opsgenie) {
	config = &client.Config{
		ApiKey: conf.ApiKey,
		Backoff: func(min, max time.Duration, attemptNum int, resp *http.Response) time.Duration {
			atomic.AddInt32(&retries, 1)
			return time.Millisecond * 1
		},
	}
}

func GetConfig() *client.Config {
	return config
}

func GetTeam() error {
	teamClient, err := team.NewClient(config)
	if err != nil {
		return err
	}
	team, err := teamClient.Get(context.TODO(), &team.GetTeamRequest{
		IdentifierType:  team.Id,
		IdentifierValue: "111-222-33",
	})
	if err != nil {
		return err
	}
	fmt.Println(team.Name)
	return nil
}
