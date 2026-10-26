# telegram-send-message

Send a message to Telegram.

## Install (master)

```
go install github.com/sikalabs/go-scripts/telegram-send-message@master
```

## Example Usage

```bash
telegram-send-message "Hello World"
```

## Credentials

The Telegram bot token and chat ID are encrypted. The encryption password
is read from the first available of:

- `SIKALABS_ENCRYPTED_PASSWORD_SIKALABS_NOTIFICATION_V1` environment variable
- `.SIKALABS_ENCRYPTED_PASSWORD_SIKALABS_NOTIFICATION_V1` (current directory)
- `~/.SIKALABS_ENCRYPTED_PASSWORD_SIKALABS_NOTIFICATION_V1`
- `/etc/SIKALABS_ENCRYPTED_PASSWORD_SIKALABS_NOTIFICATION_V1`
