package schedule

import (
	"fmt"
	"strings"
)

// OneshotCleanupTail returns the self-removal block a one-shot job's script
// ends with. It disables the timer/service, removes the files, and reloads
// systemd. Ported from tajobs jobs.go:501-506.
func OneshotCleanupTail(timerPath, servicePath string) string {
	return fmt.Sprintf(`
# self-remove the systemd unit pair once done (one-shot, not recurring)
systemctl --user disable --now "${JOB_NAME}.timer" 2>/dev/null || true
rm -f %q %q
systemctl --user daemon-reload
`, timerPath, servicePath)
}

// CycleRescheduleTail returns the self-reschedule block a custom-cycle job's
// script ends with. It advances to the next day-interval in the cycle
// (wrapping around) and rewrites its own OnCalendar= to that concrete date.
// Placeholders are substituted via strings.NewReplacer rather than
// fmt.Sprintf, since the bash itself is full of literal '%' (date format,
// modulo, printf) that would otherwise need doubling as verb-escapes.
// Ported from tajobs jobs.go:514-525.
func CycleRescheduleTail(recurPath, timerPath, timerName string) string {
	return strings.NewReplacer(
		"__RECUR_FILE__", recurPath,
		"__TIMER_PATH__", timerPath,
		"__TIMER_NAME__", timerName,
	).Replace(`
# self-reschedule for the next point in the day-interval cycle (recurring, not one-shot)
RECUR_FILE="__RECUR_FILE__"
read -r -a CYCLE < <(sed -n 1p "$RECUR_FILE")
IDX=$(sed -n 2p "$RECUR_FILE")
NEXT_IDX=$(( (IDX + 1) % ${#CYCLE[@]} ))
NEXT_DATE=$(date -d "+${CYCLE[$IDX]} days" '+%Y-%m-%d %H:%M:00')
sed -i "s/^OnCalendar=.*/OnCalendar=${NEXT_DATE}/" "__TIMER_PATH__"
printf '%s\n%d\n' "${CYCLE[*]}" "$NEXT_IDX" > "$RECUR_FILE"
systemctl --user daemon-reload
systemctl --user enable --now "__TIMER_NAME__"
`)
}
