package main

import (
    "fmt"
    "flag"
    "log"
    "bufio"
    "os"
    _ "github.com/go-sql-driver/mysql"
    "database/sql"

)

type Message struct {
    id int
    channel string
    user string
    message string
    timestamp int
}
type Channel struct {
    name string
    messages []Message
}

func ReadChannelFile(fileName string) []string {
    file, err := os.Open(fileName)
    if err != nil {
        log.Fatal(err)
    }
    defer file.Close()

    list := make([]string, 0)
    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        list = append(list, scanner.Text())
    }

    if err := scanner.Err(); err != nil {
        log.Fatal(err)
    }
    return list

}

func GetMessages(channel string, dbName string, user string) []Message {

    result := make([]Message, 0)
    conn, _ := sql.Open("mysql", user+ ":@/" + dbName)
	defer conn.Close()

	// select
    rows, err := conn.Query("SELECT * FROM message WHERE channel=\"#"+channel+ "\"")
	if err != nil {
		panic(err)
	}

    for rows.Next() {
        var id int
        var channel string
        var user string
        var message string
        var time int
        err = rows.Scan(&id, &channel, &user, &message, &time)
        if err != nil {
            log.Fatal(err)
        }
        result = append(result, Message{id,channel,user,message,time})
    }
	return result
}

func GetDbData(db string, user string, channels []string) []Channel {
    data := make([]Channel, 0)

    for _,name := range channels {
        data = append(data, Channel{
            name,
            make([]Message, 0),
            })
    }

    for _,channel := range data {
        channel.messages = GetMessages(channel.name, db, user)
        fmt.Printf("%v,%v\n", channel.name, len(channel.messages))
    }

    return data
}

func main() {
    cmdFilename := flag.String("o", "output-file", "Usage: -o <output-filename>")
    cmdDb := flag.String("db", "go-twitch-bot", "Usage: -db <dbname>")
    cmdUser := flag.String("u", "bot", "Usage: -u <dbuser>")
    cmdChannelFile := flag.String("c", "channels", "Usage: -c <channel-file>")

	flag.Parse()

	output := *cmdFilename
    db := *cmdDb
    user := *cmdUser
    channelFile := *cmdChannelFile

    println(output)
    println(db)
    println(user)
    println(channelFile)

    channels := ReadChannelFile(channelFile)

    for i, value := range channels {
        fmt.Printf("%v, %v\n", i, value)
    }

    // Pull a lot of data
    _ = GetDbData(db, user, channels)
}