# Overall

- [ ] Add logs

# Entry Page

- [ ] Make an entry page where you enter a name and color

# Room

- [ ] Implement state management via websocket messages
  - Keep state sync and chat messages as separate concerns
  - **`state_sync`** — full room state snapshot, sent once on join (or on reconnect/resync request)
    - Room ID
    - Full list of online users (name, color, icon, etc.)
  - **`state_delta`** — small event for a single change, sent after initial sync
    - `user_joined` / `user_left` / `user_updated`
    - Alternative: use separate top-level message types instead of a wrapper with an `action` field — either works
  - Client maintains local state: `state_sync` replaces it, `state_delta` patches it
  - Client can request a fresh `state_sync` if it suspects it's out of sync
  - Only send room-level state to clients that are in that room (server tracks what page each client is on)
  - [ ] Global state
  - [ ] Page-specific state (room)
- [ ] Separate chat messages from state messages — chat is an ordered event stream, state is idempotent
- [ ] Add voice chat using Pion WebRTC
  - **Do this after text chat, state management, and rooms are solid**
  - Use existing websocket as the signaling server (exchange SDP offers/answers and ICE candidates)
  - NAT traversal:
    - [ ] STUN server — helps clients discover their public IP (free ones exist like Google's)
    - [ ] TURN server — relays traffic when direct connection fails (need to host one, costs bandwidth). Without it, some users won't be able to connect
  - Phased approach:
    - [ ] Phase 1: 2 users, peer-to-peer, audio only using Pion
    - [ ] Phase 2: Get P2P working reliably before moving on
    - [ ] Phase 3: Multi-user via SFU (Selective Forwarding Unit) — server receives and forwards audio. This is basically a rewrite of the audio layer, plan for it but don't build it on day one
- [ ] Make a shared Go package for message types so server and client import the same definitions (move this up in priority — diverging types cause bugs)
- [ ] Usernames on other clients should display in that user's chosen color
- [ ] Add persistent chat by saving messages on the backend via a database
- [ ] Add a talking indicator (e.g. name turns green or similar)
