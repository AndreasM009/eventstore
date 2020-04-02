usage() {
    cat << EOF
    usage ${0} [OPTIONS]

    --namespace             Namespace name where service and secret resides.
    --input                 Input directory where to find deployment templates
    --output                Output directory to place created deployment files
EOF
exit 1
}

while [[ $# -gt 0 ]]; do 
    case ${1} in
        --namespace) namespace="$2" shift ;;
        --input) input="$2" shift ;;
        --output) output="$2" shift ;;
        *) usage ;;
    esac
    shift
done

[ -z ${namespace} ] && usage
[ -z ${input} ] && usage
[ -z ${output} ] && usage

webhookTemplate=${input}/injector-webhook-cfg-template.yaml
serviceTemplate=${input}/injector-service-template.yaml
deploymentTemplate=${input}/injector-deployment-template.yaml

resultWebhook=${output}/injector-webhook-cfg.yaml
resultService=${output}/injector-service.yaml
resultDeployment=${output}/injector-deployment.yaml
resultSecret=${output}/injector-secret.yaml

service=eventstore-injector
secret=eventstore-injector-certs

# Generate the CA cert and private key
openssl req -nodes -new -x509 -keyout ${output}/ca.key -out ${output}/ca.crt -subj "/CN=Admission Controller ${service} CA"
# Generate the private key for the webhook server
openssl genrsa -out ${output}/eventstore-injector-tls.key 2048
# Generate a Certificate Signing Request (CSR) for the private key, and sign it with the private key of the CA.
openssl req -new -key ${output}/${service}-tls.key -subj "/CN=${service}.${namespace}.svc" \
    | openssl x509 -req -CA ${output}/ca.crt -CAkey ${output}/ca.key -CAcreateserial -out ${output}/${service}-tls.crt

# create the secret with CA cert and server cert/key
kubectl create secret generic ${secret} -n ${namespace} \
        --from-file=key.key=${output}/${service}-tls.key \
        --from-file=cert.crt=${output}/${service}-tls.crt \
        --dry-run -o yaml > ${resultSecret}
    #kubectl -n ${namespace} apply -f -

ca_pem_b64="$(openssl base64 -A <"${output}/ca.crt")"

cat ${webhookTemplate} | sed -e "s|\${CA_BUNDLE}|${ca_pem_b64}|g" | sed -e "s|\${NAMESPACE}|${namespace}|g"  > ${resultWebhook}
cat ${serviceTemplate} | sed -e "s|\${NAMESPACE}|${namespace}|g" > ${resultService}
cat ${deploymentTemplate} | sed -e "s|\${NAMESPACE}|${namespace}|g" > ${resultDeployment}