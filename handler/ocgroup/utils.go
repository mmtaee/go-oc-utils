package ocgroup

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// ParseConfFile parse ocserv group config file in OcGroupConfig data type
func ParseConfFile(filename string) (*OcservGroupConfig, error) {
	var config OcservGroupConfig
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer func(file *os.File) {
		err = file.Close()
		if err != nil {
			log.Println(err)
		}
	}(file)
	scanner := bufio.NewScanner(file)
	var dnsList []string
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) == 0 || line[0] == '#' {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		switch key {
		case "rx-data-per-sec":
			config.RxDataPerSec = &value
		case "tx-data-per-sec":
			config.TxDataPerSec = &value
		case "max-same-clients":
			if val, err := strconv.Atoi(value); err == nil {
				config.MaxSameClients = &val
			}
		case "ipv4-network":
			config.IPv4Network = &value
		case "dns":
			dnsList = append(dnsList, value)
		case "no-udp":
			if val, err := strconv.ParseBool(value); err == nil {
				config.NoUDP = &val
			}
		case "keepalive":
			if val, err := strconv.Atoi(value); err == nil {
				config.KeepAlive = &val
			}
		case "dpd":
			if val, err := strconv.Atoi(value); err == nil {
				config.DPD = &val
			}
		case "mobile-dpd":
			if val, err := strconv.Atoi(value); err == nil {
				config.MobileDPD = &val
			}
		case "tunnel-all-dns":
			if val, err := strconv.ParseBool(value); err == nil {
				config.TunnelAllDNS = &val
			}
		case "restrict-user-to-routes":
			if val, err := strconv.ParseBool(value); err == nil {
				config.RestrictUserToRoutes = &val
			}
		case "stats-report-time":
			if val, err := strconv.Atoi(value); err == nil {
				config.StatsReportTime = &val
			}
		case "mtu":
			if val, err := strconv.Atoi(value); err == nil {
				config.MTU = &val
			}
		case "idle-timeout":
			if val, err := strconv.Atoi(value); err == nil {
				config.IdleTimeout = &val
			}
		case "mobile-idle-timeout":
			if val, err := strconv.Atoi(value); err == nil {
				config.MobileIdleTimeout = &val
			}
		case "session-timeout":
			if val, err := strconv.Atoi(value); err == nil {
				config.SessionTimeout = &val
			}
		}
	}
	config.DNS = &dnsList
	if err = scanner.Err(); err != nil {
		return &config, err
	}
	return &config, nil
}

// WithContext run group commands with context
func WithContext(c context.Context, operation func() error) error {
	done := make(chan error, 1)

	go func() {
		defer close(done)
		done <- operation()
	}()

	select {
	case <-c.Done():
		return fmt.Errorf("operation canceled or timed out: %w", c.Err())
	case err := <-done:
		return err
	}
}

// Writer a method to write configs in group config file
func Writer(file *os.File, config *map[string]interface{}) error {
	for k, v := range *config {
		if v == nil {
			continue
		}
		if k == "dns" {
			for _, dns := range v.([]interface{}) {
				if _, err := file.WriteString(fmt.Sprintf("dns=%s\n", dns)); err != nil {
					return fmt.Errorf("failed to write to file: %w", err)
				}
			}
			continue
		} else {
			if _, err := file.WriteString(fmt.Sprintf("%s=%v\n", k, v)); err != nil {
				return fmt.Errorf("failed to write to file: %w", err)
			}
		}
	}
	return nil
}

// OcctlExec occtl execute command
func OcctlExec(c context.Context, command string) ([]byte, error) {
	cmd := exec.CommandContext(c, "sh", "-c", fmt.Sprintf("/usr/bin/occtl %s", command))
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	return output, nil
}
