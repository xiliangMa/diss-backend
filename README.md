## run

    set env:
        MARIADB_DATABASE=diss;MARIADB_USER=diss;MARIADB_PASSWORD=abc123;MARIADB_HOST=127.0.0.1
    
    go build:
        go build -o diss-backend
        
    run：
        bee run -gendoc=true -downdoc=true
        
    run by diss-backend：
        ./diss-backend

## build image

    docker build -t dis-backend-1.0.0 .


## run in kubernetes cluster

    kubectl apply -f deploy-k8s.yml
