package main

import (
	"bufio"
	"fmt"
	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
	"net"
	"os"
	"os/exec"
	"strings"
	"time"
)

var musicCh = make(chan string, 20)
var done = make(chan bool)
var isSkip = make(chan bool)
var links = "D:/goWork/SocketServer/main/links.txt"

func main() {

	port := ":28923"
	Addr, err1 := net.ResolveTCPAddr("tcp", port)
	if err1 != nil {
		fmt.Println("Dial error5:", err1)
		panic(err1)
	}
	listener, err := net.ListenTCP("tcp", Addr)
	if err != nil {
		fmt.Println("Dial error1:", err)
		panic(err)
	}
	fmt.Println("クライアントからの入力待ち")
	if _, err := os.Stat(links); os.IsNotExist(err) {
		_, _ = os.OpenFile(links, os.O_CREATE|os.O_APPEND, 0666)
	}
	go PlayMusicLoop()
	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		fmt.Println("Accept client")
		go ListenClient(conn)
	}
}

func ListenClient(conn net.Conn) {
	defer conn.Close()
	var line string
	var messageBuf []byte
	var message string
	for {
		messageBuf = make([]byte, 1024)
		messageLen, err := conn.Read(messageBuf)
		if err != nil {
			fmt.Println("Dial error2:", err)
			panic(err)
		}
		line = string(messageBuf[:messageLen])
		line1 := strings.Replace(line, "\n", "", -1)
		line1 = strings.Replace(line1, "\r", "", -1)
		fmt.Println(line1)

		if strings.HasPrefix(line1, "/bye") {
			fmt.Println("Disconnection client")
			return
		}

		if strings.HasPrefix(line1, "/play") {
			var url = strings.Split(line1, " ")[1]
			fmt.Println(url)

			message = string(messageBuf[:messageLen])
			_, _ = conn.Write([]byte(message))

			path := getPath(url)
			if path != "" {
				musicCh <- strings.Split(path, "\\.")[0] + ".mp3"
			} else {
				fmt.Println("now downloading...")
				var cmdPy = exec.Command("python", "C:/Users/626ca/PycharmProjects/tubedoeloader/download.py ", url, "5s")
				out, err := cmdPy.Output()
				if err != nil {
					fmt.Println("Dial error6:", err)
					fmt.Println(string(out))
				}
				fmt.Println("download completed")

				path = getPath(url)

				var paths = strings.Split(path, ".")[0] + ".mp3"
				musicCh <- paths
			}
		}

		if strings.HasPrefix(line1, "/skip") {
			message = string(messageBuf[:messageLen])
			_, _ = conn.Write([]byte(message))
			isSkip <- true
		}

		if strings.HasPrefix(line1, "/list") {
			if musicCh != nil {
				for range musicCh {
					fmt.Println()
				}
			} else {
				fmt.Println("musicCh is null.")
			}
		}
	}
}

func getPath(url string) (str string) {
	var line, err1 = os.OpenFile(links, os.O_RDWR, 0666)
	if err1 != nil {
		fmt.Println("Dial error3:", err1)
		panic(err1)
	}
	defer line.Close()
	scanner := bufio.NewScanner(line)
	if err2 := scanner.Err(); err2 != nil {
		fmt.Println("Dial error4:", err2)
		panic(err2)
	}
	var path string
	for scanner.Scan() {
		tmp := strings.Split(scanner.Text(), ",")
		if tmp[0] == url {
			fmt.Println("hit.")
			path = tmp[1]
			path = "D:/goWork/SocketServer/MusicHolder/" + strings.Split(path, ".")[0]
			return path
		}
	}
	return ""
}

func PlayMusicLoop() {

	for {
		path := <-musicCh
		PlayMusicPyWrapper(strings.Split(path, ".")[0] + ".mp3")
	}
}

func PlayMusicPyWrapper(path string) {

	mp3file, _ := os.Open(path)
	s, format, err := mp3.Decode(mp3file)
	if err != nil {
		fmt.Println("Dial error7:", err)
		panic(err)
	}
	err = speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
	if err != nil {
		fmt.Println("Dial error8:", err)
		panic(err)
	}
	speaker.Play(beep.Seq(s, beep.Callback(func() {
		done <- true
	})))
	<-done

	/*out, err := exec.Command("python", "C:/Users/626ca/PycharmProjects/music_player/play_music.py", path, "5s").Output()
	if err != nil {
		fmt.Println(err)
		fmt.Println(string(out))
	}*/
}
