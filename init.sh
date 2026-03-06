#!/bin/bash

# extract wg config file path, or use default
conf="$(jq -r .config_file_path db/server/global_settings.json || echo /etc/wireguard/wg0.conf)"
iface="$(basename "${conf%.conf}")"

# manage wireguard stop/start with the container
case $WGUI_MANAGE_START in (1|t|T|true|True|TRUE)
    ip link show dev "$iface" >/dev/null 2>&1 && ip link delete dev "$iface"
    wg-quick up "$conf"
    case ${WGUI_DOWN_ON_STOP:-true} in (1|t|T|true|True|TRUE)
        trap 'wg-quick down "$conf"' SIGTERM # catches container stop
    esac
esac

# manage wireguard restarts
case $WGUI_MANAGE_RESTART in (1|t|T|true|True|TRUE)
    [[ -f $conf ]] || touch "$conf" # inotifyd needs file to exist
    inotifyd - "$conf":w | while read -r event file; do
        case ${WGUI_DOWN_ON_STOP:-true} in (1|t|T|true|True|TRUE)
            wg-quick down "$file"
        esac
        ip link show dev "$iface" >/dev/null 2>&1 && ip link delete dev "$iface"
        wg-quick up "$file"
    done &
esac


./wg-ui &
wait $!
