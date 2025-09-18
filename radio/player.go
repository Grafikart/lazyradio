package radio

import (
	"bufio"
	"context"
	"fmt"
	"net"
	"os"
	"os/exec"
	"path"
	"regexp"
	"strconv"
)

const (
	Loading = iota
	Stopped = iota
	Paused  = iota
	Playing = iota
)

type TrackItem struct {
	Name   string
	Artist string
}

func (i TrackItem) Title() string       { return i.Name }
func (i TrackItem) Description() string { return i.Artist }
func (i TrackItem) FilterValue() string { return i.Name + " " + i.Artist }

type PlayerMsg interface{}
type PlayerProgressMsg PlayerInfo
type PlayerStateChangedMsg int
type PlayerEndedMsg interface{}
type PlayerOutputMsg string
type PlayerErrorMsg error

type Player struct {
	cmd        *exec.Cmd
	ctx        context.Context
	cancel     context.CancelFunc
	ch         chan PlayerMsg
	state      int
	info       PlayerInfo
	socketFile string
}

type PlayerInfo struct {
	Duration string
	Current  string
	Progress int
}

func NewPlayer() *Player {
	ctx, cancel := context.WithCancel(context.Background())
	return &Player{
		ctx:    ctx,
		cancel: cancel,
		ch:     make(chan PlayerMsg, 5),
		state:  Stopped,
	}
}

func (p *Player) Play(s string) {
	if p.state != Stopped {
		_ = p.Stop()
	}

	p.socketFile = path.Join(os.TempDir(), "mpvsocket")
	p.ctx, p.cancel = context.WithCancel(context.Background())

	// Create the mpv command
	p.cmd = exec.CommandContext(p.ctx, "mpv",
		fmt.Sprintf(`ytdl://ytsearch1:%s`, s),
		"--ytdl-format=bestaudio",
		"--no-video",
		fmt.Sprintf("--input-ipc-server=%s", p.socketFile),
	)

	stdout, err := p.cmd.StdoutPipe()
	if err != nil {
		p.ch <- PlayerErrorMsg(fmt.Errorf("error creating stdout pipe: %v", err))
		return
	}

	progressRegex := regexp.MustCompile(`A:\s+\d{2}:(\d{2}:\d{2})\s+/\s+\d{2}:(\d{2}:\d{2})\s+\((\d+)%\)`)

	// Start the command
	if err := p.cmd.Start(); err != nil {
		p.ch <- PlayerErrorMsg(fmt.Errorf("Error starting command: %v\n", err))
		return
	}

	p.setState(Loading)

	// Read the Progress and send it to the channel
	go func() {
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			line := scanner.Text()
			matches := progressRegex.FindStringSubmatch(line)
			if len(matches) > 3 {
				if p.state == Loading {
					p.setState(Playing)
				}
				progress, _ := strconv.Atoi(matches[3])
				if progress != p.info.Progress && progress < 100 {
					p.info = PlayerInfo{
						Current:  matches[1],
						Duration: matches[2],
						Progress: progress,
					}
					select {
					case p.ch <- PlayerProgressMsg(p.info):
					default:
					}
				}
			} else {
				p.ch <- PlayerOutputMsg(line)
			}
		}
	}()

	// Wait for command completion in a goroutine
	go func() {
		defer os.Remove(p.socketFile)
		err := p.cmd.Wait()
		if err != nil {
			p.setState(Stopped)
			return
		}
		p.info = PlayerInfo{
			Progress: 100,
			Current:  p.info.Duration,
			Duration: p.info.Duration,
		}
		p.ch <- PlayerProgressMsg(p.info)
	}()
}

func (p *Player) Pause() error {
	var err error
	var newState int
	if p.state == Playing {
		err = p.sendSocket(`{ "command": ["set_property", "pause", true] }`)
		newState = Paused
	} else if p.state == Paused {
		err = p.sendSocket(`{ "command": ["set_property", "pause", false] }`)
		newState = Playing
	} else {
		return nil
	}
	if err != nil {
		p.ch <- PlayerErrorMsg(err)
		return err
	}
	p.setState(newState)
	return nil
}

func (p *Player) Stop() error {
	if p.state == Stopped {
		err := fmt.Errorf("Player already stopped")
		p.ch <- PlayerErrorMsg(err)
		return err
	}
	p.setState(Stopped)
	p.cancel()
	return nil
}

func (p *Player) setState(state int) {
	if p.state == state {
		return
	}
	p.state = state
	p.ch <- PlayerStateChangedMsg(state)
}

func (p *Player) State() int {
	return p.state
}

func (p *Player) Ch() chan PlayerMsg {
	return p.ch
}

func (p *Player) Info() PlayerInfo {
	return p.info
}

func (p *Player) sendSocket(command string) error {
	conn, err := net.Dial("unix", p.socketFile)
	if err != nil {
		return fmt.Errorf("Error connecting to socket: %v\n", err)
	}
	defer conn.Close()
	_, err = conn.Write([]byte(command + "\n"))
	if err != nil {
		return fmt.Errorf("Error writing to socket: %v\n", err)
	}
	reader := bufio.NewReader(conn)
	_, err = reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("Failed to read response: %v\n", err)
	}
	return nil
}
