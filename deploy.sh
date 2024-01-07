set -eEuo pipefail

export ACCOUNT=asafacihangir@gmail.com
export PROJ=bortecin-app
export APP=uruc-api
export TAG="gcr.io/$PROJ/$APP"
export GRUN_MIN_INSTANCE=0
export GRUN_MAX_INSTANCE=2

gcloud config set account "$ACCOUNT"
gcloud config set project "${PROJ}"


docker build -t "$TAG" .
docker push "$TAG"

gcloud beta run deploy "$APP" \
  --image "$TAG" \
  --platform "managed" \
  --region "us-central1" \
  --project "$PROJ" \
  --memory=1Gi \
  --min-instances=${GRUN_MIN_INSTANCE} \
  --max-instances=${GRUN_MAX_INSTANCE} \
  --timeout=240 \
  --no-use-http2 \
  --quiet  \
  --allow-unauthenticated