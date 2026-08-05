import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';

import 'blocs/session_cubit.dart';
import 'screens/title_screen.dart';
import 'services/api_service.dart';

Future<void> main() async {
  WidgetsFlutterBinding.ensureInitialized();
  final api = ApiService(baseUrl: 'http://192.168.1.98:8080');
  await api.init();
  runApp(DetectiveGameApp(api: api));
}

class DetectiveGameApp extends StatelessWidget {
  final ApiService api;

  const DetectiveGameApp({super.key, required this.api});

  @override
  Widget build(BuildContext context) {
    return RepositoryProvider.value(
      value: api,
      child: BlocProvider(
        create: (_) => SessionCubit(api),
        child: MaterialApp(
          title: 'ДетектИИв',
          debugShowCheckedModeBanner: false,
          theme: ThemeData(
            colorScheme: ColorScheme.fromSeed(seedColor: Colors.amber, brightness: Brightness.dark),
            useMaterial3: true,
          ),
          home: const TitleScreen(),
        ),
      ),
    );
  }
}
