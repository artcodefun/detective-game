import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';

import 'blocs/session_cubit.dart';
import 'screens/title_screen.dart';

void main() {
  runApp(const DetectiveGameApp());
}

class DetectiveGameApp extends StatelessWidget {
  const DetectiveGameApp({super.key});

  @override
  Widget build(BuildContext context) {
    return BlocProvider(
      create: (_) => SessionCubit(),
      child: MaterialApp(
        title: 'ДетектИИв',
        debugShowCheckedModeBanner: false,
        theme: ThemeData(
          colorScheme: ColorScheme.fromSeed(seedColor: Colors.amber, brightness: Brightness.dark),
          useMaterial3: true,
        ),
        home: const TitleScreen(),
      ),
    );
  }
}
