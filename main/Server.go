package main

import (
	"fmt"
	"net"
	"os"
	"os/exec"
	"strings"
)

var musicCh = make(chan string)
var links = "C:/Users/626ca/IdeaProjects/SocketChat/src/spea55/links.txt"

func main() {

	port := ":28923"
	Addr, err1 := net.ResolveTCPAddr("tcp", port)
	if err1 != nil {
		fmt.Println("Dial error5:", err1)
	}
	listener, err := net.ListenTCP("tcp", Addr)
	if err != nil {
		fmt.Println("Dial error1:", err)
	}
	fmt.Println("クライアントからの入力待ち")
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
	messageBuf := make([]byte, 1024)
	messageLen, err := conn.Read(messageBuf)
	if err != nil {
		fmt.Println("Dial error2:", err)
	}
	line = string(messageBuf[:messageLen])
	for {
		_, _ = conn.Write(messageBuf)
		fmt.Println(line)
		if strings.HasPrefix(line, "/bye") {
			return
		}

		if strings.HasPrefix(line, "/play") {
			var url = strings.SplitAfter(line, " ")[1]
			fmt.Println(url)

			message := string(messageBuf[:messageLen])

			_, _ = conn.Write([]byte(message))

			path := getPath(url)
			if &path != nil {
				musicCh <- strings.SplitAfter(path, "\\.")[0] + ".mp3"
				/*if musicCh <-strings.SplitAfter(path,"\\.")[0]+".mp3" {
					_, _ = conn.Write([]byte("Added of Queue."))
				}else {
					_, _ = conn.Write([]byte("Queue is full."))
				} */
			} else {
				var p = exec.Command("python " + "C:/Users/626ca/PycharmProjects/tubedoeloader/download.py " + url)
				fmt.Println("now downloading...")
				_ = p.Wait()
				fmt.Println("download completed")

				path = getPath(url)
				cmd := "ffmpeg -i \"" + path + "\" \"" + strings.SplitAfter(path, "\\.")[0] + ".mp3\""
				p2 := exec.Command(cmd)
				_ = p2.Wait()
				paths := strings.Split(path, "\\.")[0] + ".mp3"
				musicCh <- paths
				/* if musicCh <- paths {
					_, _ = conn.Write([]byte("Added of Queue."))
				}else {
					_, _ = conn.Write([]byte("Queue is full."))
				} */
			}
		}
	}
}

func getPath(url string) (str string) {
	var line, err1 = os.Open(links)
	if err1 != nil {
		fmt.Println("Dial error3:", err1)
	}
	var lines, err2 = line.Read([]byte(links))
	if err2 != nil {
		fmt.Println("Dial error4:", err2)
	}
	var path string
	defer line.Close()
	for {
		tmp := strings.Split(string(lines), "&,")[0]
		if tmp == url {
			fmt.Println("hit.")
			path = strings.Split(links, ",")[1]
			path = "C:/Users/626ca/IdeaProjects/SocketChat/src/spea55/" + path
			return path
		}
	}
	return path
}

func PlayMusicLoop() {
	for {
		path := <-musicCh
		_, _ = os.OpenFile(links, os.O_RDWR|os.O_CREATE, 0666)
		PlayMusicPyWrapper(strings.Split(path, "\\.")[0] + ".mp3")
	}
}

//PlayMusicPyWrapper
func PlayMusicPyWrapper(path string) {
	p := exec.Command("python C:/Users/626ca/PycharmProjects/music_player/play_music.py " + path)
	_ = p.Wait()
}
