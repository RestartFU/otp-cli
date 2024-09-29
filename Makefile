install:
	sudo rm -rf /usr/bin/otp
	go build -o otp .
	sudo mv otp /usr/bin/otp