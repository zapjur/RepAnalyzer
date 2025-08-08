#!/bin/sh
set -eu

MINIO_ALIAS="${MINIO_ALIAS:-local}"
MINIO_URL="${MINIO_URL:-http://minio:9000}"
MINIO_ROOT_USER="${MINIO_ROOT_USER:-admin}"
MINIO_ROOT_PASSWORD="${MINIO_ROOT_PASSWORD:-admin123}"
BUCKET="${BUCKET:-videos}"

ACCESS_USER="${ACCESS_USER:-access_svc}"
ACCESS_SECRET="${ACCESS_SECRET:-ACCESSSECRET}"
VIDEO_USER="${VIDEO_USER:-video_svc}"
VIDEO_SECRET="${VIDEO_SECRET:-VIDEOSECRET}"
ANALYZE_USER="${ANALYZE_USER:-analyze_svc}"
ANALYZE_SECRET="${ANALYZE_SECRET:-ANALYZESECRET}"

echo "$(date '+%F %T') Waiting for MinIO at ${MINIO_URL}..."

while :
do
  mc alias set "${MINIO_ALIAS}" "${MINIO_URL}" "${MINIO_ROOT_USER}" "${MINIO_ROOT_PASSWORD}" >/dev/null 2>&1 || true
  if mc admin info "${MINIO_ALIAS}" >/dev/null 2>&1; then
    break
  fi
  sleep 2
done
echo "$(date '+%F %T') MinIO is ready."

mc mb "${MINIO_ALIAS}/${BUCKET}" >/dev/null 2>&1 || true
mc anonymous set none "${MINIO_ALIAS}/${BUCKET}" >/dev/null 2>&1 || true

cat > /tmp/cors.json <<'JSON'
[
  {
    "AllowedOrigin": ["*"],
    "AllowedMethod": ["GET"],
    "AllowedHeader": ["*"],
    "ExposeHeader": ["ETag"],
    "MaxAgeSeconds": 3000
  }
]
JSON
mc cors set "${MINIO_ALIAS}/${BUCKET}" /tmp/cors.json >/dev/null 2>&1 || true

cat > /tmp/policy-access-read.json <<'JSON'
{ "Version":"2012-10-17", "Statement":[
  {"Effect":"Allow","Action":["s3:GetObject"],"Resource":["arn:aws:s3:::videos/*"]}
]}
JSON

cat > /tmp/policy-videos-rw.json <<'JSON'
{ "Version":"2012-10-17", "Statement":[
  {"Effect":"Allow","Action":["s3:PutObject","s3:GetObject"],"Resource":["arn:aws:s3:::videos/*"]}
]}
JSON

cat > /tmp/policy-analyze-rw.json <<'JSON'
{
  "Version":"2012-10-17",
  "Statement":[
    {
      "Effect":"Allow",
      "Action":["s3:GetObject","s3:PutObject"],
      "Resource":["arn:aws:s3:::videos/*"]
    },
    {
      "Effect":"Allow",
      "Action":["s3:ListBucket","s3:GetBucketLocation"],
      "Resource":["arn:aws:s3:::videos"]
    }
  ]
}
JSON

mc admin policy create "${MINIO_ALIAS}" access-read /tmp/policy-access-read.json >/dev/null 2>&1 || true
mc admin policy create "${MINIO_ALIAS}" videos-rw   /tmp/policy-videos-rw.json   >/dev/null 2>&1 || true
mc admin policy create "${MINIO_ALIAS}" analyze-rw  /tmp/policy-analyze-rw.json  >/dev/null 2>&1 || true

mc admin user add "${MINIO_ALIAS}" "${ACCESS_USER}"  "${ACCESS_SECRET}"  >/dev/null 2>&1 || true
mc admin user add "${MINIO_ALIAS}" "${VIDEO_USER}"   "${VIDEO_SECRET}"   >/dev/null 2>&1 || true
mc admin user add "${MINIO_ALIAS}" "${ANALYZE_USER}" "${ANALYZE_SECRET}" >/dev/null 2>&1 || true

mc admin policy attach "${MINIO_ALIAS}" access-read --user "${ACCESS_USER}"  >/dev/null 2>&1 || true
mc admin policy attach "${MINIO_ALIAS}" videos-rw   --user "${VIDEO_USER}"   >/dev/null 2>&1 || true
mc admin policy attach "${MINIO_ALIAS}" analyze-rw  --user "${ANALYZE_USER}" >/dev/null 2>&1 || true

echo "$(date '+%F %T') MinIO bootstrap done."
