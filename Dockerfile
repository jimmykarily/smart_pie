FROM scratch
ADD smart_pie /smart_pie/smart_pie
ADD public /smart_pie/public
ADD views /smart_pie/views

WORKDIR /smart_pie

CMD ["./smart_pie"]
