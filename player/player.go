package player

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"time"

	character "github.com/bwilliams-un/mud/character"
)

// Player list
var playerCount int
var players map[int]*Player

// Player state constants
const (
	StateUnknown      = 0
	StateConnecting   = 1
	StateLogin        = 2
	StateAuthenticate = 3
	StateLoading      = 4
	StateSpawn        = 5
	StateReconnect    = 6
	StatePlaying      = 10
	StateLinkdead     = 20
	StateDisconnect   = 50
)

// Player object, independent of the character
type Player struct {
	id         int
	state      byte
	lastActive time.Time
	connection net.Conn
	character  *character.Character
}

func addPlayer(player *Player) {
	if players == nil {
		players = make(map[int]*Player)
	}
	players[player.id] = player
}

func removePlayer(player *Player) {
	if player != nil {
		delete(players, player.id)
	}
}

// ConnectPlayer returns a new player object
func ConnectPlayer(conn net.Conn) {
	playerCount++
	player := Player{
		id:         playerCount,
		lastActive: time.Now(),
		connection: conn,
	}

	addPlayer(&player)
	player.SetState(StateConnecting)

	go player.handleConnection()
}

// SetState handles switching between player states
func (player *Player) SetState(state byte) {
	switch state {

	case StateConnecting:
		fmt.Printf("[Connection] Player %d connected from %s\n", player.id, player.connection.RemoteAddr())

	case StatePlaying:
		switch player.state {
		case StateLinkdead:
			fmt.Printf("[Connection] Player %d reconnected\n", player.id)
		default:
			fmt.Printf("[Connection] Player %d connected from %s\n", player.id, player.connection.RemoteAddr())
		}

	case StateLinkdead:
		fmt.Printf("[Connection] Player %d disconnected (link dead)\n", player.id)
		player.connection.Close()
		player.connection = nil
		player.lastActive = time.Now()

	}

	player.state = state
}

// handleConnection reads and processes the connection
func (player *Player) handleConnection() {
	reader := bufio.NewReader(player.connection)
	for {
		// Nanny func
		switch player.state {
		case StateConnecting:
			player.Send("Welcome!\n\n")
			player.SetState(StateLogin)

		case StateLogin:
			player.Send("What is your name? ")
			name, err := player.readInput(reader)
			if err != nil {
				player.Disconnect()
				return
			}
			player.Send("Greetings %s!\n", name)
			player.SetState(StateAuthenticate)

		case StateAuthenticate:
			player.Send("Make your mark: ")
			_, err := player.readInput(reader)
			if err != nil {
				player.Disconnect()
				return
			}

			// Check if the player is linkdead and have this player take control

			player.SetState(StateSpawn)

		case StateSpawn:
			player.Send("\nEntering the world...\n")
			player.SetState(StatePlaying)

		case StatePlaying:
			player.Send("\033[37;1mThe Void\n\033[37;0mYou stand in an empty void. You have no form and darkness permeates all around you. In the distance a faint light is your only reference.\n\n")
			player.Send("> ")
			command, err := player.readInput(reader)
			if err != nil {
				continue
			}
			player.handleCommand(&command)

		}
	}
}

func (player *Player) readInput(reader *bufio.Reader) ([]byte, error) {
	bytes, err := reader.ReadBytes('\n')
	if err != nil {
		if err != io.EOF {
			fmt.Println("[Connection] Error on read", err)
		} else {
			if player.state == StatePlaying {
				player.SetState(StateLinkdead)
			}
		}
		return nil, err
	}
	return bytes[:len(bytes)-1], nil
}

func (player *Player) handleCommand(bytes *[]byte) {

}

// Send a message to a player connection
func (player *Player) Send(format string, a ...interface{}) {
	line := fmt.Sprintf(format, a...)
	_, err := player.connection.Write([]byte(line))
	if err != nil {
		fmt.Printf("[Connection] Write failed to Player %d\n", player.id)
	}
}

// Disconnect a player
func (player *Player) Disconnect() {
	if player.connection != nil {
		player.connection.Close()
		player.connection = nil
	}
	removePlayer(player)
	fmt.Printf("[Connection] Player %d disconnected\n", player.id)
}
