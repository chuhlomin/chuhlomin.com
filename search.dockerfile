FROM chuhlomin/search:v0.1.2

ENV INDEX_PATH=/index
ENV BIND=0.0.0.0:8081

ADD search_index /index
