package commands

import (
	"testing"
	"time"

	"custom-redis/store"
)

func TestProcessSet(t *testing.T) {
	const key = "process-set-test-key"
	setTestEntry(t, key, store.Entry{
		Value:     "old value",
		ExpiresAt: time.Now().Add(time.Hour),
	})

	got, err := ProcessSet(CommandRequest{Key: key, Value: "new value"})
	assertResponse(t, got, err, "OK")

	entry := store.GetEntry(key)
	if entry.Value != "new value" {
		t.Errorf("stored value = %q, want %q", entry.Value, "new value")
	}
	if !entry.ExpiresAt.IsZero() {
		t.Errorf("stored expiration = %v, want zero time", entry.ExpiresAt)
	}
}

func TestProcessGet(t *testing.T) {
	const key = "process-get-test-key"
	setTestEntry(t, key, store.Entry{Value: "hello"})

	got, err := ProcessGet(CommandRequest{Key: key})
	assertResponse(t, got, err, "hello")

	got, err = ProcessGet(CommandRequest{Key: "process-get-missing-test-key"})
	assertResponse(t, got, err, "(nil)")
}

func TestProcessGetTreatsExpiredEntryAsMissing(t *testing.T) {
	const key = "process-get-expired-test-key"
	setTestEntry(t, key, store.Entry{
		Value:     "expired value",
		ExpiresAt: time.Now().Add(-time.Second),
	})

	got, err := ProcessGet(CommandRequest{Key: key})
	assertResponse(t, got, err, "(nil)")
}

func TestProcessDelete(t *testing.T) {
	const key = "process-delete-test-key"
	setTestEntry(t, key, store.Entry{Value: "value"})

	got, err := ProcessDelete(CommandRequest{Key: key})
	assertResponse(t, got, err, "1")
	if store.ExistsEntry(key) {
		t.Errorf("key %q still exists after delete", key)
	}

	got, err = ProcessDelete(CommandRequest{Key: key})
	assertResponse(t, got, err, "0")
}

func TestProcessExists(t *testing.T) {
	const key = "process-exists-test-key"
	setTestEntry(t, key, store.Entry{Value: "value"})

	got, err := ProcessExists(CommandRequest{Key: key})
	assertResponse(t, got, err, "1")

	store.DeleteEntry(key)
	got, err = ProcessExists(CommandRequest{Key: key})
	assertResponse(t, got, err, "0")
}

func TestProcessExistsTreatsExpiredEntryAsMissing(t *testing.T) {
	const key = "process-exists-expired-test-key"
	setTestEntry(t, key, store.Entry{
		Value:     "expired value",
		ExpiresAt: time.Now().Add(-time.Second),
	})

	got, err := ProcessExists(CommandRequest{Key: key})
	assertResponse(t, got, err, "0")
	if store.ExistsEntry(key) {
		t.Errorf("expired key %q still exists", key)
	}
}

func TestProcessExpire(t *testing.T) {
	t.Run("sets expiration for existing key", func(t *testing.T) {
		const key = "process-expire-existing-test-key"
		setTestEntry(t, key, store.Entry{Value: "value"})

		got, err := ProcessExpire(CommandRequest{Key: key, Value: "60"})
		assertResponse(t, got, err, "1")

		entry := store.GetEntry(key)
		if entry.Value != "value" {
			t.Errorf("stored value = %q, want %q", entry.Value, "value")
		}
		if entry.ExpiresAt.IsZero() || !entry.ExpiresAt.After(time.Now()) {
			t.Errorf("expiration = %v, want a future time", entry.ExpiresAt)
		}
	})

	t.Run("deletes key when expiration is zero", func(t *testing.T) {
		assertExpireDeletesKey(t, "0")
	})

	t.Run("deletes key when expiration is negative", func(t *testing.T) {
		assertExpireDeletesKey(t, "-1")
	})

	t.Run("returns zero for missing key", func(t *testing.T) {
		got, err := ProcessExpire(CommandRequest{
			Key:   "process-expire-missing-test-key",
			Value: "60",
		})
		assertResponse(t, got, err, "0")
	})

	t.Run("rejects invalid expiration", func(t *testing.T) {
		const key = "process-expire-invalid-test-key"
		setTestEntry(t, key, store.Entry{Value: "value"})

		_, err := ProcessExpire(CommandRequest{Key: key, Value: "1.5"})
		if err == nil {
			t.Fatal("ProcessExpire() error = nil, want an error")
		}

		if entry := store.GetEntry(key); entry.Value != "value" || !entry.ExpiresAt.IsZero() {
			t.Errorf("entry changed after invalid expiration: %+v", entry)
		}
	})
}

func assertExpireDeletesKey(t *testing.T, seconds string) {
	t.Helper()

	key := "process-expire-delete-test-key-" + seconds
	setTestEntry(t, key, store.Entry{Value: "value"})

	got, err := ProcessExpire(CommandRequest{Key: key, Value: seconds})
	assertResponse(t, got, err, "1")

	if store.ExistsEntry(key) {
		t.Errorf("key %q still exists after EXPIRE %s", key, seconds)
	}
}

func TestProcessTTL(t *testing.T) {
	t.Run("missing key returns minus two", func(t *testing.T) {
		got, err := ProcessTTL(CommandRequest{Key: "process-ttl-missing-test-key"})
		assertResponse(t, got, err, "-2")
	})

	t.Run("key without expiration returns minus one", func(t *testing.T) {
		const key = "process-ttl-no-expiration-test-key"
		setTestEntry(t, key, store.Entry{Value: "value"})

		got, err := ProcessTTL(CommandRequest{Key: key})
		assertResponse(t, got, err, "-1")
	})

	t.Run("empty value is treated as missing", func(t *testing.T) {
		const key = "process-ttl-empty-value-test-key"
		setTestEntry(t, key, store.Entry{})

		got, err := ProcessTTL(CommandRequest{Key: key})
		assertResponse(t, got, err, "-2")
	})

	t.Run("returns remaining whole seconds", func(t *testing.T) {
		const key = "process-ttl-future-test-key"
		setTestEntry(t, key, store.Entry{
			Value:     "value",
			ExpiresAt: time.Now().Add(2 * time.Second),
		})

		got, err := ProcessTTL(CommandRequest{Key: key})
		if err != nil {
			t.Fatalf("ProcessTTL() error = %v, want nil", err)
		}
		if got != "1" && got != "2" {
			t.Errorf("ProcessTTL() = %q, want remaining seconds 1 or 2", got)
		}
	})

	t.Run("less than one second returns minus two", func(t *testing.T) {
		const key = "process-ttl-subsecond-test-key"
		setTestEntry(t, key, store.Entry{
			Value:     "value",
			ExpiresAt: time.Now().Add(500 * time.Millisecond),
		})

		got, err := ProcessTTL(CommandRequest{Key: key})
		assertResponse(t, got, err, "-2")
		if !store.ExistsEntry(key) {
			t.Errorf("key %q was deleted before its expiration", key)
		}
	})

	t.Run("expired key returns minus two and is removed", func(t *testing.T) {
		const key = "process-ttl-expired-test-key"
		setTestEntry(t, key, store.Entry{
			Value:     "value",
			ExpiresAt: time.Now().Add(-time.Second),
		})

		got, err := ProcessTTL(CommandRequest{Key: key})
		assertResponse(t, got, err, "-2")
		if store.ExistsEntry(key) {
			t.Errorf("expired key %q still exists", key)
		}
	})
}

func TestProcessIncrement(t *testing.T) {
	testCounter(t, "INCR", ProcessIncrement, []counterTestCase{
		{name: "missing key starts at zero", want: "1"},
		{name: "increments existing value", initial: "41", want: "42"},
		{name: "rejects decimal", initial: "1.5", wantError: "value is not an integer"},
		{
			name:      "rejects out of range",
			initial:   "9223372036854775808",
			wantError: "value is not an integer or out of range",
		},
		{name: "rejects overflow", initial: "9223372036854775807", wantError: "integer overflow"},
	})
}

func TestProcessDecrement(t *testing.T) {
	testCounter(t, "DECR", ProcessDecrement, []counterTestCase{
		{name: "missing key starts at zero", want: "-1"},
		{name: "decrements existing value", initial: "41", want: "40"},
		{name: "rejects decimal", initial: "1.5", wantError: "value is not an integer"},
		{
			name:      "rejects out of range",
			initial:   "-9223372036854775809",
			wantError: "value is not an integer or out of range",
		},
		{name: "rejects overflow", initial: "-9223372036854775808", wantError: "integer overflow"},
	})
}

type counterTestCase struct {
	name      string
	initial   string
	want      string
	wantError string
}

func testCounter(
	t *testing.T,
	command string,
	process func(CommandRequest) (string, error),
	tests []counterTestCase,
) {
	t.Helper()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			key := command + "-test-" + tt.name
			if tt.initial != "" {
				setTestEntry(t, key, store.Entry{Value: tt.initial})
			} else {
				t.Cleanup(func() {
					store.DeleteEntry(key)
				})
			}

			got, err := process(CommandRequest{Key: key})
			if tt.wantError != "" {
				if err == nil || err.Error() != tt.wantError {
					t.Fatalf("error = %v, want %q", err, tt.wantError)
				}
				if stored := store.GetEntry(key).Value; stored != tt.initial {
					t.Errorf("stored value after error = %q, want unchanged value %q", stored, tt.initial)
				}
				return
			}

			assertResponse(t, got, err, tt.want)
			if stored := store.GetEntry(key).Value; stored != tt.want {
				t.Errorf("stored value = %q, want %q", stored, tt.want)
			}
		})
	}
}

func TestParseCommand(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		want      CommandRequest
		wantError string
	}{
		{
			name:  "ping",
			input: "PING",
			want:  CommandRequest{Command: CommandPing},
		},
		{
			name:  "case insensitive",
			input: "get key",
			want:  CommandRequest{Command: CommandGet, Key: "key"},
		},
		{
			name:  "set",
			input: "SET key value",
			want:  CommandRequest{Command: CommandSet, Key: "key", Value: "value"},
		},
		{
			name:  "expire",
			input: "EXPIRE key 60",
			want:  CommandRequest{Command: CommandExpire, Key: "key", Value: "60"},
		},
		{
			name:  "TTL",
			input: "TTL key",
			want:  CommandRequest{Command: CommandTTL, Key: "key"},
		},
		{name: "empty message", input: "  ", wantError: "empty message"},
		{name: "PING with arguments", input: "PING key", wantError: "PING takes no arguments"},
		{name: "GET without key", input: "GET", wantError: "GET requires one key"},
		{name: "SET without value", input: "SET key", wantError: "SET requires a key and a value"},
		{name: "EXPIRE without seconds", input: "EXPIRE key", wantError: "EXPIRE requires a key and a value"},
		{name: "EXPIRE with invalid seconds", input: "EXPIRE key 1.5", wantError: "invalid expiration time: 1.5"},
		{name: "TTL without key", input: "TTL", wantError: "TTL requires one key"},
		{name: "unknown command", input: "NOPE", wantError: "unknown command: NOPE"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseCommand(tt.input)
			if tt.wantError != "" {
				if err == nil || err.Error() != tt.wantError {
					t.Fatalf("ParseCommand() error = %v, want %q", err, tt.wantError)
				}
				return
			}
			if err != nil {
				t.Fatalf("ParseCommand() error = %v, want nil", err)
			}
			if got != tt.want {
				t.Errorf("ParseCommand() = %+v, want %+v", got, tt.want)
			}
		})
	}
}

func setTestEntry(t *testing.T, key string, entry store.Entry) {
	t.Helper()
	store.SetEntry(key, entry)
	t.Cleanup(func() {
		store.DeleteEntry(key)
	})
}

func assertResponse(t *testing.T, got string, err error, want string) {
	t.Helper()
	if err != nil {
		t.Fatalf("error = %v, want nil", err)
	}
	if got != want {
		t.Errorf("response = %q, want %q", got, want)
	}
}
