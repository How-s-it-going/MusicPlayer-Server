package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"os/exec"
	"strings"
)

var musicCh = make(chan string)
var links = "D:/go.wk/SocketServer/MusicHolder/links.txt"

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
		_, _ = os.OpenFile(links, os.O_CREATE, 0666)
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
	messageBuf := make([]byte, 1024)
	for {
		messageLen, err := conn.Read(messageBuf)
		if err != nil {
			fmt.Println("Dial error2:", err)
			panic(err)
		}
		line = string(messageBuf[:messageLen])
		fmt.Println(line)
		_, _ = conn.Write(messageBuf)
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
				var cmdPy = exec.Command("python " + "C:/Users/626ca/PycharmProjects/tubedoeloader/download.py " + url)
				fmt.Println("now downloading...")
				_ = cmdPy.Start()
				_ = cmdPy.Wait()
				fmt.Println("download completed")

				path = getPath(url)
				cmd := "ffmpeg -i \"" + path + "\" \"" + strings.SplitAfter(path, "\\.")[0] + ".mp3\""
				_ = exec.Command(cmd).Run()
				paths := strings.SplitAfter(path, "\\.")[0] + ".mp3"
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
	var line, err1 = os.OpenFile(links, os.O_RDONLY, 0666)
	if err1 != nil {
		fmt.Println("Dial error3:", err1)
		panic(err1)
	}
	scanner := bufio.NewScanner(line)
	if err2 := scanner.Err(); err2 != nil {
		fmt.Println("Dial error4:", err2)
		panic(err2)
	}
	var path string
	defer line.Close()
	for scanner.Scan() {
		tmp := strings.SplitAfter(scanner.Text(), "&,")[0]
		if tmp == url {
			fmt.Println("hit.")
			path = strings.SplitAfter(links, ",")[1]
			path = "D:/go.wk/SocketServer/MusicHolder/" + path
			return path
		}
	}
	path = ""
	return path
}

func PlayMusicLoop() {
	for {
		path := <-musicCh
		PlayMusicPyWrapper(strings.Split(path, "\\.")[0] + ".mp3")
	}
}

//PlayMusicPyWrapper
func PlayMusicPyWrapper(path string) {
	_ = exec.Command("python C:/Users/626ca/PycharmProjects/music_player/play_music.py " + path).Run()
}
