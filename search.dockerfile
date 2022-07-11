FROM ghcr.io/chuhlomin/search

ENV INDEX_PATH=/index
ENV BIND=0.0.0.0:8081

ADD search_index /index
