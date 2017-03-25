import sublime, sublime_plugin, os
from subprocess import call


class gotemplator(sublime_plugin.EventListener):

	ext = ".gtm"

	def __init__(self):
		settings = sublime.load_settings('gotemplator.sublime-settings')
		gotemplator.ext = settings.get("extension")

	def on_post_save(self, view):
		filename = view.file_name()
		if filename is None:
			return

		if filename.endswith(gotemplator.ext) is not True:
			return

		folder = os.path.dirname(view.file_name())
		call(['gotemplator', '-e', gotemplator.ext, folder])
