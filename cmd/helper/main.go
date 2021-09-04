package main

import (
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/gateway"
	statepkg "github.com/diamondburned/arikawa/v3/state"
	_ "github.com/joho/godotenv/autoload"
	"github.com/termora/berry/commands/static/rpc"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"
)

var (
	log *zap.SugaredLogger

	token   = os.Getenv("TOKEN")
	guildID = discord.GuildID(mustSnowflake(os.Getenv("GUILD_ID")))
	rpcHost = os.Getenv("RPC")

	state  *statepkg.State
	client rpc.GuildMemberServiceClient

	userMu sync.Mutex
	users  = map[discord.UserID]discord.Member{}
)

func main() {
	initLog()

	conn, err := grpc.Dial(rpcHost, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("Could not connect to RPC server: %v", err)
	}
	client = rpc.NewGuildMemberServiceClient(conn)
	log.Infof("Connected to RPC server")

	state, err = statepkg.NewWithIntents("Bot "+os.Getenv("TOKEN"), gateway.IntentGuildMembers|gateway.IntentGuilds)
	if err != nil {
		log.Fatalf("Error opening state: %v\n", err)
	}
	log.Infof("Created state")

	state.AddHandler(guildCreate)
	state.AddHandler(guildMemberAdd)
	state.AddHandler(guildMemberRemove)
	state.AddHandler(guildMemberUpdate)
	state.AddHandler(guildMemberChunk)

	err = state.Open(state.Context())
	if err != nil {
		log.Fatalf("Error opening connection to Discord: %v", err)
	}
	log.Infof("Connected to Discord")

	defer func() {
		state.Close()
		log.Info("Disconnected from Discord")
	}()

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	log.Info("Interrupt signal received. Shutting down...")
}

func mustSnowflake(s string) discord.Snowflake {
	sf, err := discord.ParseSnowflake(s)
	if err != nil {
		panic(err)
	}
	return sf
}

func initLog() {
	// set up a logger
	zcfg := zap.NewProductionConfig()
	zcfg.Encoding = "console"
	zcfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	zcfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	zcfg.EncoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
	zcfg.EncoderConfig.EncodeDuration = zapcore.StringDurationEncoder
	zcfg.Level.SetLevel(zapcore.DebugLevel)

	zap, err := zcfg.Build(zap.AddStacktrace(zapcore.ErrorLevel))
	if err != nil {
		panic(err)
	}
	log = zap.Sugar()
}
