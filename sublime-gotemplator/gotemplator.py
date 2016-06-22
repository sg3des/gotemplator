import sublime, sublime_plugin, os

class goinstall(sublime_plugin.EventListener):

	ext = ".gtm"

	def __init__(self):
		settings = sublime.load_settings('gotemplator.sublime-settings')
		goinstall.ext = settings.get("extension")

	def on_post_save(self, view):
		if view.file_name().endswith(goinstall.ext) is not True:
			return

		folder = os.path.dirname(view.file_name())
		view.window().run_command('exec',{'cmd':['gotemplator','-e',goinstall.ext],'working_dir':folder})
