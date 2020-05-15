#!/usr/bin/env bash
echo ">make build"
make build

execute(){
    slug=$1
    shift
    commands=("$@")

    echo -e "*** $slug ***\n"
    for cmd in "${commands[@]}"
        do
            echo "> ./edgex-cli $cmd"
            ./edgex-cli $cmd
            echo -e "\n"
        done
    echo -e "*** End $slug ***\n"
}
test_deviceService() {
    declare -a commands=("deviceservice  list --no-pager"
                         )
    execute "Device Service" "${commands[@]}"
}

test_device() {
#TODO add device fails because there is no deviceservice.
    declare -a commands=("device  list --no-pager"
#                         "device  add -f samples/createDevice.toml"
#                         "device  rm --name Car-001"
                         )
    execute "Device" "${commands[@]}"
}

test_dp() {
    declare -a commands=("profile list --no-pager"
                         "profile add samples/createDP.yaml"
                         "profile rm --name Simple-Device-1")
    execute "DeviceProfile" "${commands[@]}"
}


test_intervals(){
    declare -a commands=("interval add samples/createInterval.toml"
                         "interval list --no-pager" \
                         "interval rm --name noon"\
                         "interval rm --name fourteen-hundrend-hours")
    execute "Intervals" "${commands[@]}"
}


test_notifications(){
    declare -a commands=("notification add samples/createNotifications.toml"
                         "notification list --no-pager" \
                         "notification list --slug notice-001" \
                         "notification list --labels=temperature")
    execute "Notification" "${commands[@]}"
}

test_others(){
    declare -a commands=("addressable list --no-pager"
                         "deviceservice list --no-pager" \
                         "event list --no-pager" \
                         "reading list --no-pager"\
                         "subscription list --no-pager"
                         "status"
                         "version")
    execute "Others" "${commands[@]}"
}

test_all(){
      test_deviceService
      test_dp
      test_device
      test_notifications
      test_intervals
      test_others
}

if [[ $# -eq 0 ]] ; then
    test_all
    exit 0
fi

for i in "$@"
do
    case $i in
     -d)
        test_device
        ;;
     -dp)
      test_dp
      ;;
     -ds)
      test_deviceService
      ;;
     -i)
      test_intervals
      ;;
     -n)
      test_notifications
      ;;
     -o)
      test_others
      ;;

      *)
      test_all
      ;;

    esac
done