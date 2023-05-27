package scan

import (
	"context"
	"fmt"
	"net"
	"time"
)

func Run(hl *HostsList, ports []int, udp, open, closed bool, ttl time.Duration) ([]Results, error) {
	res := make([]Results, 0, len(hl.Hosts))

	if ttl == 0 {
		ttl = 1 * time.Hour
	}

	ctx, cancel := context.WithTimeout(context.Background(), ttl)
	defer cancel()

	resCh := make(chan []Results)

	go func() {
		for _, h := range hl.Hosts {
			r := Results{
				Host: h,
			}

			if _, err := net.LookupHost(h); err != nil {
				r.NotFound = true
				res = append(res, r)
				continue
			}

			for _, p := range ports {
				state := scanPort(h, p, udp)
				if open && !bool(state.Open) {
					continue
				}

				if closed && bool(state.Open) {
					continue
				}

				r.PortStates = append(r.PortStates, state)
			}

			res = append(res, r)
		}

		resCh <- res
		close(resCh)
	}()

	for {
		select {
		case res := <-resCh:
			return res, nil
		case <-ctx.Done():
			err := ctx.Err()
			return nil, err
		}
	}
}

type PortState struct {
	Port int
	Open state
}

type state bool

func (s state) String() string {
	if s {
		return "open"
	}

	return "closed"
}

func scanPort(host string, port int, udp bool) PortState {
	p := PortState{
		Port: port,
	}

	address := net.JoinHostPort(host, fmt.Sprintf("%d", port))

	network := "tcp"
	if udp {
		network = "udp"
	}

	scanConn, err := net.DialTimeout(network, address, 1*time.Second)
	if err != nil {
		return p
	}

	scanConn.Close()
	p.Open = true
	return p
}

type Results struct {
	Host       string
	NotFound   bool
	PortStates []PortState
}
