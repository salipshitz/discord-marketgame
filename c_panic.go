package main

import (
	"errors"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func cmdPanic(s *discordgo.Session, m *discordgo.Message, args []string) error {
	if len(args) < 1 {
		return errors.New("Too few arguments.\nUsage: !panic <symbol | *>")
	}
	
	u, err := GetUser(m.Author.ID)
	if err != nil {
		return err
	}
	
	if args[0] == "*" || args[0] == "." {
		for symbol := range u.Shares {
			err = cmdPanic(s, m, []string{symbol});
			if err != nil {
				displayError(s, m, e)
			}
		}
		
		return nil
	} else {
		symbol := strings.ToUpper(args[0])
				
		shares, ok := u.Shares[symbol]
		if !ok {
			return errors.New("You don't own any " + symbol + ".")
		}

		return cmdSell(s, m, []string{symbol, strconv.Itoa(shares)})
	}
}

