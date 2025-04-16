.PHONY: all run-all test clean

all:
	@echo "Available commands:"
	@echo "  make run-all   # Run all servers concurrently and then execute tests"
	@echo "  make test      # Run test suite only"
	@echo "  make clean     # Clean background processes/logs (if needed)"


run-all:
	@echo "Starting servers on ports 9001, 9002, and 9003..."
	# Start servers in background, output logs to files.
	nohup go run server/main.go 9001 > server9001.log 2>&1 &
	nohup go run server/main.go 9002 > server9002.log 2>&1 &
	nohup go run server/main.go 9003 > server9003.log 2>&1 &
	@echo "Servers started."
	@sleep 2
	@echo "Running client...Client will run the test cases."
	go run client/main.go


test:
	go run client/main.go


clean:
	rm -f server9001.log server9002.log server9003.log
	@echo "Cleaned log files."
