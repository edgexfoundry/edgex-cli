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
#                         "device  add -f samples/createDevice.json"
#                         "device  rm --name Car-001"
                         )
    execute "Device" "${commands[@]}"
}

test_dp() {
    declare -a commands=("profile list --no-pager"
                         "profile add -f samples/createDP.json"
                         "profile rm --name DeviceProfileCLI")
    execute "DeviceProfile" "${commands[@]}"
}


test_intervals(){
    declare -a commands=("interval add -f samples/createInterval.json"
                         "interval list --no-pager" \
                         "interval rm --name noon"\
                         "interval rm --name fourteen-hundrend-hours")
    execute "Intervals" "${commands[@]}"
}


test_notifications(){
    declare -a commands=("notification add -f samples/createNotification.json"
                         "notification list --new --no-pager" \
                         "notification list --sender SystemManagement" \
                         "notification list --slug notice-001" \
                         "notification list --labels=temperature" \
                         "notification rm --slug notice-001")
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