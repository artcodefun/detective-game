import 'dart:async';

import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:provider/provider.dart';

import 'screens/title_screen.dart';
import 'services/audio_service.dart';
import 'services/api_service.dart';
import 'services/session_service.dart';

Future<void> main() async {
  WidgetsFlutterBinding.ensureInitialized();
  await SystemChrome.setPreferredOrientations([DeviceOrientation.portraitUp]);
  const apiBaseUrl = String.fromEnvironment('API_BASE_URL', defaultValue: 'http://localhost:8080');
  final api = ApiService(baseUrl: apiBaseUrl);
  final audio = AudioService();
  await audio.initialize();
  unawaited(api.initialize());
  unawaited(audio.resumeMusic());
  runApp(DetectiveGameApp(api: api, audio: audio));
}

class DetectiveGameApp extends StatefulWidget {
  final ApiService api;
  final AudioService audio;

  DetectiveGameApp({super.key, required this.api, AudioService? audio}) : audio = audio ?? AudioService();

  @override
  State<DetectiveGameApp> createState() => _DetectiveGameAppState();
}

class _DetectiveGameAppState extends State<DetectiveGameApp> with WidgetsBindingObserver {
  @override
  void initState() {
    super.initState();
    WidgetsBinding.instance.addObserver(this);
  }

  @override
  void dispose() {
    WidgetsBinding.instance.removeObserver(this);
    widget.audio.pauseMusic();
    super.dispose();
  }

  @override
  void didChangeAppLifecycleState(AppLifecycleState state) {
    switch (state) {
      case AppLifecycleState.resumed:
        unawaited(widget.audio.resumeAfterAppLifecycle());
        return;
      case AppLifecycleState.inactive:
      case AppLifecycleState.paused:
      case AppLifecycleState.detached:
      case AppLifecycleState.hidden:
        widget.audio.pauseForAppLifecycle();
    }
  }

  @override
  Widget build(BuildContext context) {
    return MultiProvider(
      providers: [
        ChangeNotifierProvider.value(value: widget.api),
        Provider.value(value: widget.audio),
        ChangeNotifierProvider(create: (_) => SessionService(widget.api)),
      ],
      child: MaterialApp(
        title: 'ДетектИИв',
        debugShowCheckedModeBanner: false,
        // Android renders Flutter apps edge-to-edge. Keep every route and
        // modal sheet above the system navigation area; AppBar still handles
        // the top system inset itself.
        builder: (context, child) => SafeArea(top: false, child: child ?? const SizedBox.shrink()),
        theme: ThemeData(
          colorScheme: ColorScheme.fromSeed(seedColor: Colors.amber, brightness: Brightness.dark),
          useMaterial3: true,
        ),
        home: const TitleScreen(),
      ),
    );
  }
}
