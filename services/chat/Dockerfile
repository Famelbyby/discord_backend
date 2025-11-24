FROM python:3.12

WORKDIR /chat

COPY ./ /chat


EXPOSE 9091
RUN pip install poetry==1.8.2 
RUN poetry install --no-root

