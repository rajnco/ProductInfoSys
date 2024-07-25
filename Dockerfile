FROM golang:1.22.5-bookworm

RUN apt-get update && \
	apt-get upgrade -y  && \
	apt-get clean && \ 
	apt-get autoclean

WORKDIR /app

COPY . .

RUN groupadd -r nonroot && useradd -g nonroot nonroot

RUN chown -R nonroot:nonroot /app

RUN mkdir -p /home/nonroot/.cache

RUN chown -R nonroot:nonroot /home/nonroot

USER nonroot

RUN go mod download


RUN go build -o /home/nonroot/productinfosys

EXPOSE 8080

CMD [ "/home/nonroot/productinfosys" ]
