MINIO_URL = http://localhost:9000
MINIO_ALIAS = local
MINIO_ROOT_USER = admin
MINIO_ROOT_PASS = admin123

BUCKET = videos
USER_PREFIX = auth0_67d1b4eec471877e05db4b83
EXERCISE_NAME = Deadlift
WEIGHT ?=
REPS ?=

ORIG_SRC     ?= ./test-files/deadlift.mp4
BARPATH_SRC  ?= ./test-files/deadlift-barpath.mp4
BARPATH_CSV  ?= ./test-files/deadlift-barpath.csv
BARPATH_META ?= ./test-files/deadlift-barpath_meta.json
POSE_SRC     ?= ./test-files/deadlift-pose.mp4
POSE_CSV     ?= ./test-files/deadlift-pose.csv
POSE_META 	 ?= ./test-files/deadlift-pose_meta.json

ORIG_DST     = $(BUCKET)/$(USER_PREFIX)/Deadlift/original/deadlift.mp4
BARPATH_DST  = $(BUCKET)/$(USER_PREFIX)/Deadlift/barpath/deadlift.mp4
BARPATH_CSV_DST  = $(BUCKET)/$(USER_PREFIX)/Deadlift/barpath/deadlift.csv
BARPATH_META_DST  = $(BUCKET)/$(USER_PREFIX)/Deadlift/barpath/deadlift_meta.json
POSE_DST     = $(BUCKET)/$(USER_PREFIX)/Deadlift/pose/deadlift.mp4
POSE_CSV_DST     = $(BUCKET)/$(USER_PREFIX)/Deadlift/pose/deadlift.csv
POSE_META_DST  = $(BUCKET)/$(USER_PREFIX)/Deadlift/pose/deadlift_meta.json

PGHOST ?= localhost
PGPORT ?= 5432
PGDATABASE ?= database
PGUSER ?= user
PGPASSWORD ?= password
USER_ID ?= 1

all: upload-videos db-insert

alias-root:
	@echo "Setting mc alias as root..."
	@mc alias set $(MINIO_ALIAS) $(MINIO_URL) $(MINIO_ROOT_USER) $(MINIO_ROOT_PASS) >/dev/null

ensure-bucket: alias-root
	@echo "Ensuring bucket $(BUCKET) exists..."
	@mc ls $(MINIO_ALIAS)/$(BUCKET) >/dev/null 2>&1 || mc mb $(MINIO_ALIAS)/$(BUCKET)

upload-videos: ensure-bucket
	@echo "Uploading original video as root..."
	@mc cp $(ORIG_SRC)    $(MINIO_ALIAS)/$(ORIG_DST)
	@echo "Uploading barpath video as root..."
	@mc cp $(BARPATH_SRC) $(MINIO_ALIAS)/$(BARPATH_DST)
	@echo "Uploading barpath csv as root..."
	@mc cp $(BARPATH_CSV) $(MINIO_ALIAS)/$(BARPATH_CSV_DST)
	@echo "Uploading barpath metadata as root..."
	@mc cp $(BARPATH_META) $(MINIO_ALIAS)/$(BARPATH_META_DST)
	@echo "Uploading pose-estimation video as root..."
	@mc cp $(POSE_SRC)    $(MINIO_ALIAS)/$(POSE_DST)
	@echo "Uploading pose-estimation csv as root..."
	@mc cp $(POSE_CSV)    $(MINIO_ALIAS)/$(POSE_CSV_DST)
	@echo "Uploading pose-estimation metadata as root..."
	@mc cp $(POSE_META)   $(MINIO_ALIAS)/$(POSE_META_DST)
	@echo "Upload finished"

db-insert:
	@echo "Inserting database rows..."
	@ID=$$(PGHOST=$(PGHOST) PGPORT=$(PGPORT) PGDATABASE=$(PGDATABASE) PGUSER=$(PGUSER) PGPASSWORD=$(PGPASSWORD) \
	psql -v ON_ERROR_STOP=1 -qtAX --no-psqlrc -c "\
		INSERT INTO videos (user_id, bucket, object_key, weight, repetitions, exercise_name) \
		VALUES ($(USER_ID), '$(BUCKET)', '$(USER_PREFIX)/Deadlift/original/deadlift.mp4', \
		        NULLIF('$(WEIGHT)','')::INT, NULLIF('$(REPS)','')::INT, '$(EXERCISE_NAME)') \
		RETURNING id;"); \
	ID=$$(echo $$ID | tr -d '[:space:]'); \
	echo "Video id: $$ID"; \
	PGHOST=$(PGHOST) PGPORT=$(PGPORT) PGDATABASE=$(PGDATABASE) PGUSER=$(PGUSER) PGPASSWORD=$(PGPASSWORD) \
	psql -v ON_ERROR_STOP=1 --no-psqlrc -c "\
		INSERT INTO video_analysis (video_id, type, bucket, object_key, metadata) \
		VALUES ($$ID, 'barpath', '$(BUCKET)', '$(USER_PREFIX)/Deadlift/barpath/deadlift.mp4', '{}'::jsonb);" ; \
	PGHOST=$(PGHOST) PGPORT=$(PGPORT) PGDATABASE=$(PGDATABASE) PGUSER=$(PGUSER) PGPASSWORD=$(PGPASSWORD) \
	psql -v ON_ERROR_STOP=1 --no-psqlrc -c "\
		INSERT INTO video_analysis (video_id, type, bucket, object_key, metadata) \
		VALUES ($$ID, 'pose', '$(BUCKET)', '$(USER_PREFIX)/Deadlift/pose/deadlift.mp4', '{}'::jsonb);"
	@echo "Database insert finished"

