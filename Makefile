NAME:=map
build: 
	go build -o $(NAME) .


.SILENT: build
