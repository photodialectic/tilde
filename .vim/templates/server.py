import os.path
import tornado.ioloop
import tornado.web

import json
import urllib
import markdown
import base64

BASE_DIR = os.path.dirname(__file__)

class MainHandler(tornado.web.RequestHandler):
    def get(self):
        self.render('index.html', message="Hello, world")

class HealthHandler(tornado.web.RequestHandler):
    def get(self):
        self.write('healthy')

def make_app():
    debug = True
    app_settings = {
        'debug' : debug,
        'static_path': os.path.join(BASE_DIR, 'static'),
        'template_path': os.path.join(BASE_DIR),
        'static_url_prefix' : '/static/',
    }

    if(debug):
        app_settings['static_url_prefix'] = '//localhost:6121/static/'


    routes = [
        (r"/", MainHandler),
        (r"/health-check", HealthHandler),
    ]

    return tornado.web.Application(routes, **app_settings)

if __name__ == "__main__":
    app = make_app()
    app.listen(3000)
    tornado.ioloop.IOLoop.current().start()
