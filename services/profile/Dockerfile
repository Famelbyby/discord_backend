FROM python:3.12

WORKDIR /profile

COPY ./ /profile


EXPOSE 9000
RUN pip install poetry==1.8.2 
RUN poetry install --no-root

