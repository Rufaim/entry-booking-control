package main

import "os/user"

func GetUserIdentifier() string {
	user, err := user.Current()
	if err != nil {
		panic(err)
	}
	return user.Username
}

func PanicOnError(err error) {
	if err != nil {
		panic(err)
	}
}
