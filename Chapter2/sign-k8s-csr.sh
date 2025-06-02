#!/bin/bash

csr_name="my-client-csr"
name="${1:-my-user}"
csr="${2:-$csr_name.csr}"

if [[ ! -f "$csr" ]]; then
  echo "CSR file '$csr' not found!"
  exit 1
fi

echo "Sending CSR to Kubernetes..."

cat <<EOF | kubectl create -f -
apiVersion: certificates.k8s.io/v1
kind: CertificateSigningRequest
metadata:
  name: ${csr_name}
spec:
  groups:
  - system:authenticated
  request: $(cat ${csr} | base64 | tr -d '\n')
  usages:
  - key encipherment
  - client auth
EOF

echo
echo "Approving signing request..."
kubectl certificate approve "${csr_name}"

echo
echo "Downloading certificate..."
kubectl get csr "${csr_name}" -o jsonpath='{.status.certificate}' | base64 --decode > "$(basename ${csr} .csr).crt"

echo
echo "Cleaning up"
kubectl delete csr "${csr_name}"

echo
echo "✅ Add the following to the 'users' list in your kubeconfig file:"
echo "- name: ${name}"
echo "  user:"
echo "    client-certificate: ${PWD}/$(basename ${csr} .csr).crt"
echo "    client-key: ${PWD}/$(basename ${csr} .csr)-key.pem"

echo
echo "📌 Next you may want to add a role-binding for this user."
