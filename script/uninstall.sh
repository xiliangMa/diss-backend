uninstall_image() {
    if [ $# -eq 0 ]; then
        echo "Image name is required for uninstalling operation"
        return 1
    fi

    # shellcheck disable=SC2006
    images=`docker images | grep $1 | awk -F ' ' '{printf "%s:%s\n",$1,$2}'`
    test -z "${images}" && return

    for image in ${images}
    do
        if [ $1 = "scope" ]; then
            # shellcheck disable=SC2006
            tmp=`echo "$1" | grep diss`
            [ -z ${tmp} ] && return
        fi

        echo "Uninstalling image ${image}"

        # shellcheck disable=SC2006
        containers=`docker ps -a | grep ${image} | awk -F' ' '{print $NF}'`
        if [ -n "${containers}" ]; then
            for container in ${containers}
            do
                echo "Stopping container ${container}"
                if ! docker stop ${container} > /dev/null; then
                    echo "Stopping container ${container} failed"
                fi

                echo "Removing container ${container}"
                if ! docker rm -f ${container} > /dev/null; then
                    echo "Removing container ${container} failed"
                fi
            done
        fi

        echo "Removing image ${image}"
        if ! docker rmi ${image} > /dev/null; then
            echo "Uninstalling image ${image} failed"
            continue
        fi
        echo "Uninstalling image ${image} success"
    done
}


for image in  registry.cn-qingdao.aliyuncs.com/diss/ timescale-db diss-db
do
    uninstall_image ${image}
done


echo ""
echo "----------- REMOVE DATA and INSTALLPACK , exec follow step ----------"
echo "if want remove log, delete /var/log/diss-metric dir"
echo "if want remove install package , delete /opt/diss-agent-offline dir"
echo ""


