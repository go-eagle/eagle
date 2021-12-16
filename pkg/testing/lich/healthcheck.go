package lich

import (
	"database/sql"
	"fmt"
	"net"
	"strconv"
	"strings"

	"github.com/go-eagle/eagle/pkg/log"
	// Register go-sql-driver stuff
	_ "github.com/go-sql-driver/mysql"
)

var healthchecks = map[string]func(*Container) error{"mysql": checkMysql, "mariadb": checkMysql}

// Healthcheck check container health.
func (c *Container) Healthcheck() (err error) {
	if status, health := c.State.Status, c.State.Health.Status; !c.State.Running || (health != "" && health != "healthy") {
		err = fmt.Errorf("service: %s | container: %s not running", c.GetImage(), c.GetID())
		log.Error("docker status(%s) health(%s) error(%v)", status, health, err)
		return
	}
	if check, ok := healthchecks[c.GetImage()]; ok {
		err = check(c)
		return
	}
	for proto, ports := range c.NetworkSettings.Ports {
		if id := c.GetID(); !strings.Contains(proto, "tcp") {
			log.Errorf("container: %s helloworld(%s) unsupported.", id, proto)
			continue
		}
		for _, publish := range ports {
			var (
				ip      = net.ParseIP(publish.HostIP)
				port, _ = strconv.Atoi(publish.HostPort)
				tcpAddr = &net.TCPAddr{IP: ip, Port: port}
				tcpConn *net.TCPConn
			)
			if tcpConn, err = net.DialTCP("tcp", nil, tcpAddr); err != nil {
				log.Errorf("net.DialTCP(%s:%s) error(%v)", publish.HostIP, publish.HostPort, err)
				return
			}
			_ = tcpConn.Close()
		}
	}
	return
}

func checkMysql(c *Container) (err error) {
	var ip, port, user, passwd string
	for _, env := range c.Config.Env {
		splits := strings.Split(env, "=")
		if strings.Contains(splits[0], "MYSQL_ROOT_PASSWORD") {
			user, passwd = "root", splits[1]
			continue
		}
		if strings.Contains(splits[0], "MYSQL_ALLOW_EMPTY_PASSWORD") {
			user, passwd = "root", ""
			continue
		}
		if strings.Contains(splits[0], "MYSQL_USER") {
			user = splits[1]
			continue
		}
		if strings.Contains(splits[0], "MYSQL_PASSWORD") {
			passwd = splits[1]
			continue
		}
	}
	var db *sql.DB
	if ports, ok := c.NetworkSettings.Ports["3306/tcp"]; ok {
		ip, port = ports[0].HostIP, ports[0].HostPort
	}
	var dsn = fmt.Sprintf("%s:%s@tcp(%s:%s)/", user, passwd, ip, port)
	if db, err = sql.Open("mysql", dsn); err != nil {
		log.Errorf("sql.Open(mysql) dsn(%s) error(%v)", dsn, err)
		return
	}
	if err = db.Ping(); err != nil {
		log.Errorf("ping(db) dsn(%s) error(%v)", dsn, err)
	}
	defer func() {
		_ = db.Close()
	}()
	return
}
