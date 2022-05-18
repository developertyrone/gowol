package controller

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-ping/ping"
	"github.com/sabhiram/go-wol/wol"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"
	. "web/model"
)

func Router(r *gin.Engine) {
	r.GET("/", index)
	r.POST("/api/wol", wolservice)
	r.POST("/api/ping", pingservice)

}

func index(c *gin.Context) {

	jsonFile, err := os.Open("config.json")
	// if we os.Open returns an error then handle it
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Successfully Opened config.json")
	// defer the closing of our jsonFile so that we can parse it later on

	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	// we initialize our Users array
	var targets Targets

	// we unmarshal our byteArray which contains our
	// jsonFile's content into 'users' which we defined above
	json.Unmarshal(byteValue, &targets)

	//var result map[string]interface{}
	//json.Unmarshal([]byte(byteValue), &result)

	//fmt.Println(targets.Targets)

	for i := 0; i < len(targets.Targets); i++ {
		fmt.Println("Target Name: " + targets.Targets[i].Name)

	}

	c.HTML(
		http.StatusOK,
		"views/index.html",
		gin.H{
			"Targets": targets.Targets,
		},
	)
}

func wolservice(c *gin.Context) {
	var targetHost Target
	var err error
	c.BindJSON(&targetHost)
	err = wakeCmd(targetHost.Macaddr)
	if err != nil {
		log.Println(err)
		c.String(http.StatusInternalServerError, err.Error())
	}
	c.String(http.StatusOK, "Packet Sent")
}

func pingservice(c *gin.Context) {
	var targetHost Target
	var err error
	c.BindJSON(&targetHost)
	stats, err := pingCmd(targetHost.Ip)
	//err = longPingCmd(targetHost.Ip)
	if err != nil {
		log.Println(err)
		c.String(http.StatusInternalServerError, err.Error())
	}
	log.Println(err)
	if stats.PacketsRecv > 0 {
		c.String(http.StatusOK, "Active")
	} else {
		c.String(http.StatusOK, "Ping Timeout")
	}
}

// Run the wake command.
func wakeCmd(macAddr string) error {

	// The address to broadcast to is usually the default `255.255.255.255` but
	// can be overloaded by specifying an override in the CLI arguments.
	bcastAddr := fmt.Sprintf("%s:%s", "255.255.255.255", "9")
	udpAddr, err := net.ResolveUDPAddr("udp", bcastAddr)
	if err != nil {
		return err
	}

	// Build the magic packet.
	mp, err := wol.New(macAddr)
	if err != nil {
		return err
	}

	// Grab a stream of bytes to send.
	bs, err := mp.Marshal()
	if err != nil {
		return err
	}

	// Grab a UDP connection to send our packet of bytes.
	conn, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		return err
	}
	defer conn.Close()

	fmt.Printf("Attempting to send a magic packet to MAC %s\n", macAddr)
	fmt.Printf("... Broadcasting to: %s\n", bcastAddr)
	n, err := conn.Write(bs)
	if err == nil && n != 102 {
		err = fmt.Errorf("magic packet sent was %d bytes (expected 102 bytes sent)", n)
	}
	if err != nil {
		return err
	}

	fmt.Printf("Magic packet sent successfully to %s\n", macAddr)
	return nil
}

func pingCmd(Ip string) (*ping.Statistics, error) {
	pinger, err := ping.NewPinger(Ip)
	if err != nil {
		return nil, err
	}
	pinger.Count = 3
	pinger.Timeout = time.Second * 5
	err = pinger.Run() // Blocks until finished.
	if err != nil {
		return nil, err
	}

	stats := pinger.Statistics() // get send/receive/duplicate/rtt stats
	return stats, nil
}

func longPingCmd(Ip string) error {
	pinger, err := ping.NewPinger(Ip)
	if err != nil {
		panic(err)
	}

	// Listen for Ctrl-C.
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for _ = range c {
			pinger.Stop()
		}
	}()

	pinger.OnRecv = func(pkt *ping.Packet) {
		fmt.Printf("%d bytes from %s: icmp_seq=%d time=%v\n",
			pkt.Nbytes, pkt.IPAddr, pkt.Seq, pkt.Rtt)
	}

	pinger.OnDuplicateRecv = func(pkt *ping.Packet) {
		fmt.Printf("%d bytes from %s: icmp_seq=%d time=%v ttl=%v (DUP!)\n",
			pkt.Nbytes, pkt.IPAddr, pkt.Seq, pkt.Rtt, pkt.Ttl)
	}

	pinger.OnFinish = func(stats *ping.Statistics) {
		fmt.Printf("\n--- %s ping statistics ---\n", stats.Addr)
		fmt.Printf("%d packets transmitted, %d packets received, %v%% packet loss\n",
			stats.PacketsSent, stats.PacketsRecv, stats.PacketLoss)
		fmt.Printf("round-trip min/avg/max/stddev = %v/%v/%v/%v\n",
			stats.MinRtt, stats.AvgRtt, stats.MaxRtt, stats.StdDevRtt)
	}

	fmt.Printf("PING %s (%s):\n", pinger.Addr(), pinger.IPAddr())
	err = pinger.Run()
	if err != nil {
		panic(err)
	}
	return nil
}

func portConnectCmd(host string, ports []string) error {
	for _, port := range ports {
		timeout := time.Second
		conn, err := net.DialTimeout("tcp", net.JoinHostPort(host, port), timeout)
		if err != nil {
			return err
		}
		if conn != nil {
			defer conn.Close()
			return nil
		}
	}

	return nil
}
