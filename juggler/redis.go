package juggler

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/PuerkitoBio/exp/juggler/msg"
	"github.com/garyburd/redigo/redis"
)

func pullRedisRes(c *Conn) {
	const minTimeoutSecs = 1

	rc := c.srv.CallPool.Get()
	defer rc.Close()

	for {
		// check for stop signal
		select {
		case <-c.kill:
			return
		default:
		}

		// get the next call result
		toSecs := int(c.srv.ResBrpopTimeout / time.Second)
		if toSecs <= minTimeoutSecs {
			toSecs = minTimeoutSecs
		}
		b, err := redis.Bytes(rc.Do("BRPOP", fmt.Sprintf(resKey, c.UUID), toSecs))
		if err != nil {
			// TODO : do not return
		}

		var m resPayload
		if err := json.Unmarshal(b, &m); err != nil {
			// TODO
		}

		// check if it is still expected (not timed-out)
		cnt, err := rc.Do("DEL", fmt.Sprintf(resTimeoutKey, c.UUID, m.MsgUUID))
		if err != nil {
			// TODO
		}

		if cnt == 1 {
			res := msg.NewRes(m)
			c.Send(res)
		}
	}
}
