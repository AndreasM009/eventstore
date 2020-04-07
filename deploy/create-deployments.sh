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

webhookInjectorTemplate=${input}/injector-webhook-cfg-template.yaml
serviceInjectorTemplate=${input}/injector-service-template.yaml
deploymentInjectorTemplate=${input}/injector-deployment-template.yaml
serviceOperatorTemplate=${input}/operator-service-template.yaml
deploymentOperatorTemplate=${input}/operator-deployment-template.yaml
operatorServiceAccountTemplate=${input}/operator-service-account-template.yaml
operatorClusterRoleTemplate=${input}/operator-clusterrole-template.yaml
operatorClusterRoleBindingTemplate=${input}/operator-clusterrole-binding-template.yaml


resultInjectorWebhook=${output}/injector-webhook-cfg.yaml
resultInjectorService=${output}/injector-service.yaml
resultInjectorDeployment=${output}/injector-deployment.yaml
resultInjectorSecret=${output}/injector-secret.yaml
resultOperatorDeployment=${output}/operator-deployment.yaml
resultOperatorService=${output}/operator-service.yaml
resultOperatorServiceAccount=${output}/operator-service-account.yaml
resultOperatorClusterRole=${output}/operator-clusterrole.yaml
resultOperatorClusterRoleBinding=${output}/operator-clusterrole-binding.yaml

serviceInjector=eventstore-injector
secretInjector=eventstore-injector-certs

# Generate the CA cert and private key
openssl req -nodes -new -x509 -keyout ${output}/ca.key -out ${output}/ca.crt -subj "/CN=Admission Controller ${serviceInjector} CA"
# Generate the private key for the webhook server
openssl genrsa -out ${output}/eventstore-injector-tls.key 2048
# Generate a Certificate Signing Request (CSR) for the private key, and sign it with the private key of the CA.
openssl req -new -key ${output}/${serviceInjector}-tls.key -subj "/CN=${serviceInjector}.${namespace}.svc" \
    | openssl x509 -req -CA ${output}/ca.crt -CAkey ${output}/ca.key -CAcreateserial -out ${output}/${serviceInjector}-tls.crt

# create the secret with CA cert and server cert/key
kubectl create secret generic ${secretInjector} -n ${namespace} \
        --from-file=key.key=${output}/${serviceInjector}-tls.key \
        --from-file=cert.crt=${output}/${serviceInjector}-tls.crt \
        --dry-run -o yaml > ${resultInjectorSecret}
    #kubectl -n ${namespace} apply -f -

ca_pem_b64="$(openssl base64 -A <"${output}/ca.crt")"

cat ${webhookInjectorTemplate} | sed -e "s|\${CA_BUNDLE}|${ca_pem_b64}|g" | sed -e "s|\${NAMESPACE}|${namespace}|g"  > ${resultInjectorWebhook}
cat ${serviceInjectorTemplate} | sed -e "s|\${NAMESPACE}|${namespace}|g" > ${resultInjectorService}
cat ${deploymentInjectorTemplate} | sed -e "s|\${NAMESPACE}|${namespace}|g" > ${resultInjectorDeployment}
cat ${serviceOperatorTemplate} | sed -e "s|\${NAMESPACE}|${namespace}|g" > ${resultOperatorService}
cat ${deploymentOperatorTemplate} | sed -e "s|\${NAMESPACE}|${namespace}|g" > ${resultOperatorDeployment}
cat ${operatorServiceAccountTemplate} | sed -e "s|\${NAMESPACE}|${namespace}|g" > ${resultOperatorServiceAccount}
cat ${operatorClusterRoleTemplate} | sed -e "s|\${NAMESPACE}|${namespace}|g" > ${resultOperatorClusterRole}
cat ${operatorClusterRoleBindingTemplate} | sed -e "s|\${NAMESPACE}|${namespace}|g" > ${resultOperatorClusterRoleBinding}