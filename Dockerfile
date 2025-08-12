# Dockerfile
FROM python:3.10.16-bullseye

WORKDIR /tempbuild
# Copy App Over App
COPY . /tempbuild

RUN go build .

# Install Python Modules
RUN python -m pip install --upgrade pip
RUN pip install -r /srv/requirements.txt




# Launch
CMD ["python","run.py"]
