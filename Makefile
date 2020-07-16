install:
	go install github.com/OakAnderson/internetTester/cmd/nettest
	go install github.com/OakAnderson/internetTester/cmd/nettest-csv
	go install github.com/OakAnderson/internetTester/cmd/nettest-mysql
	./API/linux/speedtest

build:
	cd cmd/nettest/ && go build
	cd cmd/nettest-csv/ && go build
	cd cmd/nettest-mysql/ && go build

update:
	go install github.com/OakAnderson/internetTester/cmd/nettest
	go install github.com/OakAnderson/internetTester/cmd/nettest-csv
	go install github.com/OakAnderson/internetTester/cmd/nettest-mysql
