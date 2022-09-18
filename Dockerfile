FROM golang 

WORKDIR /app/ctp 

ENV ORG_OVERRIDE=catppuccin-rfc

COPY go.mod go.sum ./
RUN go mod download 

RUN mkdir -p /root/.vscode/extensions
RUN mkdir -p /root/.config/helix/themes
COPY . .

RUN go build -v -o ./out/ctp .

ENTRYPOINT [ "./out/ctp" ]
