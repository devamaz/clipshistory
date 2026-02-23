package app

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/devamaz/clipshistory/internal/store"
	"golang.design/x/clipboard"
)

func MonitorClipboard(storage *store.Store) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()
	check := make(chan bool)
	nonViolate := make(chan bool)
	var prevByte []byte
	var curByte []byte
	changed := clipboard.Watch(context.Background(), clipboard.FmtText)

	go func(ctx context.Context) {
		for {
			<-check
			str := string(curByte)
			fmt.Println("Read string: " + str)
			if !strings.Contains(str, "keyword here") {
				fmt.Println("content non violated")
				nonViolate <- false
			} else {
				fmt.Println("content violated")
				nonViolate <- true
			}
		}
	}(ctx)

	for {
		var buffer string = ""
		curByte = <-changed
		buffer += string(curByte)
		if strings.Compare(string(prevByte), string(curByte)) != 0 {
			check <- true
			isViolate := <-nonViolate
			if !isViolate {
				// Save clipboard content to database
				now := time.Now().Unix()
				clip := store.Clip{
					Content:      string(curByte),
					CreatedAt:    now,
					LastCopiedAt: now,
					CopyCount:    1,
					IsPinned:     false,
					IsDeleted:    false,
				}
				if err := storage.Save(clip); err != nil {
					fmt.Printf("Error saving clip to database: %v\n", err)
				} else {
					fmt.Println("Successfully saved clip to database")
				}
				fmt.Println("Going to write: " + string(curByte))
				clipboard.Write(clipboard.FmtText, curByte)
			} else {
				fmt.Println("Going to write nothing here")
				clipboard.Write(clipboard.FmtText, []byte("\n"))
				curByte = []byte("\n")
			}
		}
		prevByte = curByte
	}
}
