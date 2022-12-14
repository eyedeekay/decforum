FROM golang
COPY . /src/go/src/github.com/eyedeekay/decforum
WORKDIR /src/go/src/github.com/eyedeekay/decforum
RUN go build -o /bin/decforum
WORKDIR /var/lib/decforum
CMD ["/bin/decforum"]
