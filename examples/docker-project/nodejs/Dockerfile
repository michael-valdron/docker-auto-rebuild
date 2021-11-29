# Pull base image.
FROM node:12-alpine

# Copy project files to 
COPY ${PWD} /data

# Define working directory.
WORKDIR /data

# Install Dependencies
RUN npm install

# Expose port
EXPOSE 3000

# Entrypoint command.
ENTRYPOINT [ "npm" ]

# Define default command.
CMD [ "start" ]
