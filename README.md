# font2png

- This is a tool that converts each character from fonts into PNG files for creating datasets for machine learning.
- Place font files in the `imports` directory and run it, image files will be generated in the `exports` directory.
- Change `config.yaml` according to your needs.
- Set the range of unicode you want to generate using `unicodestart` and `unicodeend`.
- By setting `captions: true`, you can generate caption files for each image.

```shell
# Development
docker compose up -d --build

# Production
docker compose -f compose.yaml up -d --build
```
